# anthropic
Zero dependency (Unofficial) Go client for the Anthropic API.

### Installation

```
go get github.com/fabiustech/anthropic
```

### Example Usage

```go
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/fabiustech/anthropic"
)

var key *string

func init() {
	key = flag.String("key", "", "api key")
	flag.Parse()
}

func main() {
	var client = anthropic.NewClient(*key)
	var resp, err = client.NewCompletion(context.Background(), &anthropic.Request{
		Prompt:            anthropic.NewPromptFromString("Tell me a haiku about trees"),
		Model:             anthropic.Claude,
		MaxTokensToSample: 300,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Completion)
}
```

### Contributing

Contributions are welcome and encouraged! Feel free to report any bugs / feature requests as issues.
