package anthropic

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/bedrockruntime"
)

// NewBedrockClient returns a new client for the Bedrock API.
// I've implemented this client to have identical signatures to the original client.
func NewBedrockClient(p client.ConfigProvider, cfgs ...*aws.Config) *BedrockClient {
	return &BedrockClient{
		client: bedrockruntime.New(p, cfgs...),
	}
}

type BedrockClient struct {
	client *bedrockruntime.BedrockRuntime
	debug  bool
}

// Debug enables debug logging. When enabled, the client will log the request's prompt.
func (bc *BedrockClient) Debug() {
	bc.debug = true
}

// NewCompletion returns a completion response from the API.
func (bc *BedrockClient) NewCompletion(ctx context.Context, req *Request) (*Response, error) {
	var m = req.Model
	req.Model = UnknownModel

	var b, err = json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var resp *bedrockruntime.InvokeModelOutput
	resp, err = bc.client.InvokeModelWithContext(ctx, &bedrockruntime.InvokeModelInput{
		Body:    b,
		ModelId: aws.String(m.BedrockString()),
	})
	if err != nil {
		return nil, err
	}

	var out = &Response{}
	if err = json.Unmarshal(resp.Body, out); err != nil {
		return nil, err
	}

	return out, nil
}

// NewCompletionStreamedBatchResponse returns a completion response from the API, which appears to the caller
// as a non-streaming response. However, it is actually a streaming response under the hood. This is useful
// in cases where you are getting a 524 error from the API, which is caused by the API taking too long to
// respond. Our theory is that these errors are caused by the API taking too long to respond to the load balancer,
// which then closes the connection. Since a streaming request will get a response as soon as the API has
// generated the first token, this should prevent the load balancer from closing the connection.
//
// Note: This may be deprecated at any time, but is currently needed as most requests are running into this issue.
func (bc *BedrockClient) NewCompletionStreamedBatchResponse(ctx context.Context, req *Request) (*Response, error) {
	var resps, errs, err = bc.NewStreamingCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	var resp = &Response{}

loop:
	for {
		select {
		case err = <-errs:
			if err == nil {
				break loop
			}
			return nil, err
		case rr, ok := <-resps:
			if !ok {
				break loop
			}
			resp.Completion += rr.Completion
		}
	}

	return resp, nil
}

// NewStreamingCompletion returns two channels: the first will be sent |*Response|s as they are received from
// the API and the second is sent any error(s) encountered while receiving / parsing responses.
func (bc *BedrockClient) NewStreamingCompletion(ctx context.Context, req *Request) (<-chan *Response, <-chan error, error) {
	var m = req.Model
	req.Model = UnknownModel

	var b, err = json.Marshal(req)
	if err != nil {
		return nil, nil, err
	}

	var resp *bedrockruntime.InvokeModelWithResponseStreamOutput
	resp, err = bc.client.InvokeModelWithResponseStreamWithContext(ctx, &bedrockruntime.InvokeModelWithResponseStreamInput{
		Body:    b,
		ModelId: aws.String(m.BedrockString()),
	})
	if err != nil {
		return nil, nil, err
	}

	var respCh = make(chan *Response)
	var errCh = make(chan error)

	go func() {
		defer close(respCh)
		defer close(errCh)

		var s = resp.GetStream()
		var events = s.Events()
		for {
			select {
			case ev := <-events:
				switch pp := ev.(type) {
				case *bedrockruntime.PayloadPart:
					var out = &Response{}
					if err = json.Unmarshal(pp.Bytes, out); err != nil {
						errCh <- err
						return
					}

					respCh <- out

					if out.StopReason != nil {
						return
					}
				case *bedrockruntime.ResponseStreamUnknownEvent:
					// Continue.
				}

			default:
				// Continue.
			}

			if err = s.Err(); err != nil {
				errCh <- err
				break
			}
		}

	}()

	return respCh, errCh, nil
}
