# tello-sdk-go

Go WebSocket SDK for the Tello `/sdk` protocol.

```go
package main

import (
	"context"
	"log"

	"github.com/tello-ai/tello-sdk-go/tello"
)

func main() {
	ctx := context.Background()
	client, err := tello.NewClient("")
	if err != nil {
		log.Fatal(err)
	}

	client.On(tello.EventTypeUserTurn, func(ctx context.Context, event tello.Event) error {
		return client.Answer(ctx, "heard: "+event.Text, "", "")
	})
	// Send DTMF digits: client.SendDtmf(ctx, "1234#", "", "")

	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if err := client.CreateCall(ctx, "+821012345678", "agent-1", "reservation check", nil, ""); err != nil {
		log.Fatal(err)
	}
	if err := client.WaitClosed(ctx); err != nil {
		log.Fatal(err)
	}
}
```

`NewClient("")` reads `TELLO_API_KEY`; `WithURL` overrides the default `ws://localhost:3000/sdk`.

`Connect` authenticates the API key internally: after the socket opens it sends an `authenticate` frame and returns only once the server confirms with `auth.ok`. No `Authorization` header or query-string token is used, and the key never appears in logs or error messages. `Connect` returns an error if authentication fails, the server closes with code `4401`, or `auth.ok` does not arrive within `WithOpenTimeout`. No commands run before authentication completes.

The module path assumes this package is mirrored to `github.com/tello-ai/tello-sdk-go` for release. If it remains inside a monorepo, set the module path to the fetchable repository path before publishing.
