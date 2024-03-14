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
	ticker := ontick.New(context.TODO(), 1*time.Second, 4, "key42")

	go func() {
		time.Sleep(42 * time.Second)
		ticker.Stop()
	}()

	var count atomic.Int64

	ticker.Do(
		func(ctx context.Context) {
			count.Add(1)

			t := ticker.GetTickTimeFromContext(ctx)

			if t.Second()%2 != 0 {
				fmt.Printf("%v [%02d] | Look at me I'm a rock star!\n", t.Format(time.UnixDate), count.Load())
			}
		},
		func(ctx context.Context) {
			count.Add(1)

			t := ticker.GetTickTimeFromContext(ctx)

			if t.Second()%2 == 0 {
				fmt.Printf("%v [%02d] | I have a cunning plan...\n", t.Format(time.UnixDate), count.Load())
			}
		},
	)

	ticker.Do(func(ctx context.Context) {
		fmt.Println("This will not be executed.")
	})

	ticker.Wait()
}
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue to discuss potential improvements or features.

## License

`OnTick` is available under the MIT license. See the LICENSE file for more info.
