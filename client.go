package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	v3 "github.com/fabiustech/anthropic/v3"
)

const (
	host               = "api.anthropic.com"
	completionEndpoint = "v1/complete"
	messagesEndpoint   = "v1/messages"
	apiKeyHeader       = "X-Api-Key"
	apiVersionHeader   = "Anthropic-Version"
	defaultVersion     = "2023-06-01"

	// Header and value to enable using the beta version of the API which allows for a max output tokens of 8192.
	// https://docs.anthropic.com/en/release-notes/api#july-15th-2024
	betaHeaderName             = "anthropic-beta"
	betaOutputTokenHeaderValue = "max-tokens-3-5-sonnet-2024-07-15"
	betaPromptCacheHeaderValue = "prompt-caching-2024-07-31"
)

// Client is a client for the Anthropic API.
type Client struct {
	key   string
	debug bool
	// requestHeaders is a map of custom headers to be sent with each request.
	requestHeaders http.Header
}

// NewClient returns a client with the given API key.
func NewClient(key string) *Client {
	return &Client{key: key, requestHeaders: http.Header{apiVersionHeader: {defaultVersion}}}
}

// SetVersion set's the value passed in the |Anthropic-Version| header for requests.
// The default value is "2023-06-01".
func (c *Client) SetVersion(version string) {
	if c.requestHeaders == nil {
		c.requestHeaders = make(http.Header)
	}

	c.requestHeaders.Set(apiVersionHeader, version)
}

// AddRequestHeaders adds the custom headers to be sent with each request.
func (c *Client) AddRequestHeaders(headers http.Header) {
	if c.requestHeaders == nil {
		return
	}

	for k, v := range headers {
		c.requestHeaders[k] = v
	}
}

// SetBetaMaxOutputTokenHeader sets the |anthropic-beta| header to "max-tokens-3-5-sonnet-2024-07-15".
func (c *Client) SetBetaMaxOutputTokenHeader() {
	if c.requestHeaders == nil {
		c.requestHeaders = make(http.Header)
	}

	c.requestHeaders.Add(betaHeaderName, betaOutputTokenHeaderValue)
}

// SetBetaPromptCacheHeader sets the |anthropic-beta| header to "prompt-caching-2024-07-31".
func (c *Client) SetBetaPromptCacheHeader() {
	if c.requestHeaders == nil {
		c.requestHeaders = make(http.Header)
	}

	c.requestHeaders.Add(betaHeaderName, betaPromptCacheHeaderValue)
}

// Debug enables debug logging. When enabled, the client will log the request's prompt.
func (c *Client) Debug() {
	c.debug = true
}

// NewCompletion returns a completion response from the API.
func (c *Client) NewCompletion(ctx context.Context, req *Request) (*Response, error) {
	if c.debug {
		slog.Info("completion request", "prompt", req.Prompt)
	}

	var b, err = c.post(ctx, completionEndpoint, req)
	if err != nil {
		return nil, err
	}

	var resp = &Response{}
	if err = json.Unmarshal(b, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

type streamingMessageRequest[T v3.RequestMessage] struct {
	*v3.Request[T]
	Stream bool `json:"stream"`
}

type v3Event struct {
	Type         string             `json:"type"`
	Index        int                `json:"index"`
	Delta        *v3.MessageContent `json:"delta,omitempty"`
	ContentBlock *v3.MessageContent `json:"content_block,omitempty"`
}

// NewMessageRequest makes a request to the messages endpoint.
func (c *Client) NewMessageRequest(ctx context.Context, req *v3.Request[v3.Message]) (*v3.Response, error) {
	if c.debug {
		req.Debug()
	}

	var b, err = c.post(ctx, messagesEndpoint, req)
	if err != nil {
		return nil, err
	}

	var resp = &v3.Response{}
	if err = json.Unmarshal(b, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) NewStreamingMessageRequest(ctx context.Context, req *v3.Request[v3.Message]) (*v3.Response, <-chan string, <-chan error, error) {
	if c.debug {
		req.Debug()
	}

	var receive, errs, err = c.postStream(ctx, messagesEndpoint, &streamingMessageRequest[v3.Message]{
		Request: req,
		Stream:  true,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	var respCh = make(chan string)
	var errCh = make(chan error)

	var resp = &v3.Response{}

	go func() {
		defer close(respCh)
		defer close(errCh)

		for {
			select {
			case b := <-receive:
				var events []*event
				events, err = parseEvents(b)
				if err != nil {
					errCh <- err
					return
				}

				for _, e := range events {
					switch e.Type {
					case eventTypeMessageStart:
						if err = json.Unmarshal(e.Data, resp); err != nil {
							errCh <- err
							return
						}
					case eventTypeMessageDelta:
						var delta = &v3.Response{}
						if err = json.Unmarshal(e.Data, delta); err != nil {
							errCh <- err
							return
						}

						resp.StopReason = delta.StopReason
						resp.StopSequence = delta.StopSequence
						resp.Usage = delta.Usage
					case eventTypeMessageStop:
						return
					case eventTypeContentBlockStart:
						var block = &v3Event{}
						if err = json.Unmarshal(e.Data, block); err != nil {
							errCh <- err
							return
						}

						resp.Content = append(resp.Content, block.ContentBlock)
						respCh <- block.ContentBlock.Text
					case eventTypeContentBlockDelta:
						var delta = &v3Event{}
						if err = json.Unmarshal(e.Data, delta); err != nil {
							errCh <- err
							return
						}

						resp.Content[len(resp.Content)-1].Text += delta.Delta.Text
						respCh <- delta.Delta.Text
					case eventTypeContentBlockStop:
						continue
					case eventTypeError:
						var errResp = &ResponseError{}
						if err = json.Unmarshal(e.Data, errResp); err != nil {
							errCh <- errors.New(string(e.Data))
							return
						}

						errCh <- errResp
						return
					case eventTypePing:
						// Do nothing.
					default:
						errCh <- ErrBadEvent
						return
					}
				}
			case err = <-errs:
				errCh <- err
				return
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			}
		}
	}()

	return resp, respCh, errCh, nil
}

func (c *Client) NewStreamingShortHandMessageRequest(ctx context.Context, req *v3.Request[v3.ShortHandMessage]) (*v3.Response, <-chan string, <-chan error, error) {
	if c.debug {
		req.Debug()
	}

	var receive, errs, err = c.postStream(ctx, messagesEndpoint, &streamingMessageRequest[v3.ShortHandMessage]{
		Request: req,
		Stream:  true,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	var respCh = make(chan string)
	var errCh = make(chan error)

	var resp = &v3.Response{}

	go func() {
		defer close(respCh)
		defer close(errCh)

		for {
			select {
			case b := <-receive:
				var events []*event
				events, err = parseEvents(b)
				if err != nil {
					errCh <- err
					return
				}

				for _, e := range events {
					switch e.Type {
					case eventTypeMessageStart:
						if err = json.Unmarshal(e.Data, resp); err != nil {
							errCh <- err
							return
						}
					case eventTypeMessageDelta:
						var delta = &v3.Response{}
						if err = json.Unmarshal(e.Data, delta); err != nil {
							errCh <- err
							return
						}

						resp.StopReason = delta.StopReason
						resp.StopSequence = delta.StopSequence
						resp.Usage = delta.Usage
					case eventTypeMessageStop:
						return
					case eventTypeContentBlockStart:
						var block = &v3Event{}
						if err = json.Unmarshal(e.Data, block); err != nil {
							errCh <- err
							return
						}

						resp.Content = append(resp.Content, block.ContentBlock)
						respCh <- block.ContentBlock.Text
					case eventTypeContentBlockDelta:
						var delta = &v3Event{}
						if err = json.Unmarshal(e.Data, delta); err != nil {
							errCh <- err
							return
						}

						resp.Content[len(resp.Content)-1].Text += delta.Delta.Text
						respCh <- delta.Delta.Text
					case eventTypeContentBlockStop:
						continue
					case eventTypeError:
						var errResp = &ResponseError{}
						if err = json.Unmarshal(e.Data, errResp); err != nil {
							errCh <- errors.New(string(e.Data))
							return
						}

						errCh <- errResp
						return
					case eventTypePing:
						// Do nothing.
					default:
						errCh <- ErrBadEvent
						return
					}
				}
			case err = <-errs:
				errCh <- err
				return
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			}
		}
	}()

	return resp, respCh, errCh, nil
}

// NewShortHandMessageRequest makes a request to the messages endpoint.
func (c *Client) NewShortHandMessageRequest(ctx context.Context, req *v3.Request[v3.ShortHandMessage]) (*v3.Response, error) {
	if c.debug {
		req.Debug()
	}

	var b, err = c.post(ctx, messagesEndpoint, req)
	if err != nil {
		return nil, err
	}

	var resp = &v3.Response{}
	if err = json.Unmarshal(b, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// NewCompletionStreamedBatchResponse returns a completion response from the API, which appears to the caller
// as a non-streaming response. However, it is actually a streaming response under the hood. This is useful
// in cases where you are getting a 524 error from the API, which is caused by the API taking too long to
// respond. Our theory is that these errors are caused by the API taking too long to respond to the load balancer,
// which then closes the connection. Since a streaming request will get a response as soon as the API has
// generated the first token, this should prevent the load balancer from closing the connection.
//
// Note: This may be deprecated at any time, but is currently needed as most requests are running into this issue.
func (c *Client) NewCompletionStreamedBatchResponse(ctx context.Context, req *Request) (*Response, error) {
	var resps, errs, err = c.NewStreamingCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	var resp = &Response{}

	for {
		select {
		case err = <-errs:
			return nil, err
		case rr := <-resps:
			resp.Completion += rr.Completion
			if rr.StopReason != nil {
				resp.StopReason = rr.StopReason
				resp.Stop = rr.Stop
				resp.Model = rr.Model
			}
		}

		if resp.StopReason != nil {
			break
		}
	}

	return resp, nil
}

type streamingRequest struct {
	*Request
	Stream bool `json:"stream"`
}

// NewStreamingCompletion returns two channels: the first will be sent |*Response|s as they are received from
// the API and the second is sent any error(s) encountered while receiving / parsing responses.
func (c *Client) NewStreamingCompletion(ctx context.Context, req *Request) (<-chan *Response, <-chan error, error) {
	if c.debug {
		slog.Info("completion request", "prompt", req.Prompt)
	}

	var receive, errs, err = c.postStream(ctx, completionEndpoint, &streamingRequest{
		Request: req,
		Stream:  true,
	})
	if err != nil {
		return nil, nil, err
	}
	var respCh = make(chan *Response)
	var errCh = make(chan error)

	go func() {
		defer close(respCh)
		defer close(errCh)

		for {
			select {
			case b := <-receive:
				var events []*event
				events, err = parseEvents(b)
				if err != nil {
					errCh <- err
					return
				}

				for _, e := range events {
					switch e.Type {
					case eventTypeCompletion:
						var resp = &Response{}

						if err = json.Unmarshal(e.Data, resp); err != nil {
							errCh <- err
							return
						}

						respCh <- resp

						if resp.StopReason != nil {
							return
						}
					case eventTypeError:
						var errResp = &ResponseError{}
						if err = json.Unmarshal(e.Data, errResp); err != nil {
							errCh <- errors.New(string(e.Data))
							return
						}

						errCh <- errResp
						return
					case eventTypePing:
						// Do nothing.
					default:
						errCh <- ErrBadEvent
						return

					}
				}
			case err = <-errs:
				errCh <- err
				return
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			}
		}
	}()

	return respCh, errCh, nil
}

var re = regexp.MustCompile("event: (.*?)\ndata: (.*?)\n")

type eventType string

const (
	// Completion event types.
	eventTypeCompletion eventType = "completion"

	// Message event types.
	eventTypeMessageStart      eventType = "message_start"
	eventTypeMessageDelta      eventType = "message_delta"
	eventTypeMessageStop       eventType = "message_stop"
	eventTypeContentBlockStart eventType = "content_block_start"
	eventTypeContentBlockDelta eventType = "content_block_delta"
	eventTypeContentBlockStop  eventType = "content_block_stop"

	// Event types for all streaming requests.
	eventTypeError eventType = "error"
	eventTypePing  eventType = "ping"
)

// ErrBadEvent is returned when an event is received that cannot be parsed.
var ErrBadEvent = errors.New("bad event")

type event struct {
	Type eventType
	Data []byte
}

func parseEvents(b []byte) ([]*event, error) {
	var out []*event

	var matches = re.FindAllSubmatch(b, -1)
	for _, group := range matches {
		if len(group) != 3 {
			return nil, ErrBadEvent
		}

		var ev = &event{
			Type: eventType(strings.TrimSpace(string(group[1]))),
			Data: group[2],
		}
		out = append(out, ev)
	}

	return out, nil
}

func (c *Client) post(ctx context.Context, path string, payload any) ([]byte, error) {
	var b, err = json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var u = url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path,
	}

	var req *http.Request
	req, err = c.newRequest(ctx, "POST", u.String(), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = interpretResponse(resp); err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

const bufferSize = 1024

func (c *Client) postStream(ctx context.Context, path string, payload any) (<-chan []byte, <-chan error, error) {
	var b, err = json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	var u = url.URL{
		Scheme: "https",
		Host:   host,
		Path:   path,
	}

	var req *http.Request
	req, err = c.newRequest(ctx, "POST", u.String(), bytes.NewBuffer(b))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "text/event-stream; charset=utf-8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "no-cache")

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if err = interpretResponse(resp); err != nil {
		_ = resp.Body.Close()
		return nil, nil, err
	}

	var events = make(chan []byte)
	var errCh = make(chan error)

	go func() {
		defer resp.Body.Close()
		defer close(events)
		defer close(errCh)

		for {
			var msg = make([]byte, bufferSize)
			_, err = resp.Body.Read(msg)

			switch {
			case errors.Is(err, io.EOF):
				return
			case err != nil:
				errCh <- err
				return
			case ctx.Err() != nil:
				errCh <- ctx.Err()
				return
			default:
				// No-op.
			}

			events <- msg
		}
	}()

	return events, errCh, nil
}

func (c *Client) newRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	var req, err = http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(apiKeyHeader, c.key)

	if c.requestHeaders != nil {
		for k, v := range c.requestHeaders {
			req.Header[k] = v
		}
	}

	return req, nil
}

func interpretResponse(resp *http.Response) error {
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var b, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("code: %d, unable to read response body", resp.StatusCode)

		}

		var errResp = &ResponseError{}
		if err = json.Unmarshal(b, errResp); err == nil {
			errResp.Err.Code = resp.StatusCode
			return errResp
		}

		return fmt.Errorf("code: %d, error: %s", resp.StatusCode, string(b))
	}

	return nil
}
