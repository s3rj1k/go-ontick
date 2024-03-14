# OnTick

`OnTick` is a Go module that provides the ability to execute function(s) on timer tick. 

## Getting Started

### Prerequisites

- Go 1.22.1 or later.

### Installation

To install `OnTick`, use the `go get` command:

```bash
go get github.com/s3rj1k/go-ontick
```

### Example

```
package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/s3rj1k/go-ontick"
)

func main() {
	var (
		c1 atomic.Int64
		c2 atomic.Int64
	)

	ticker := ontick.New(context.TODO(), 10*time.Millisecond, 2, "key42")

	go func() {
		<-time.After(42 * 10 * time.Millisecond)
		ticker.Stop()
	}()

	ticker.Do(
		func(ctx context.Context) {
			c1.Add(1)
		},
		func(ctx context.Context) {
			c2.Add(1)
		},
	)

	ticker.Do(func(ctx context.Context) {
		panic("This will not be executed.")
	})

	ticker.Wait()

	fmt.Printf("%d =? %d\n", c1.Load(), c2.Load())
}
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue to discuss potential improvements or features.

## License

`OnTick` is available under the MIT license. See the LICENSE file for more info.
