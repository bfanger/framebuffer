A golang library providing low-level access to the framebuffer.

Inspired by
https://github.com/xdsopl/framebuffer
https://github.com/gonutz/framebuffer
https://github.com/kaey/framebuffer

But instead of providing an abstraction, this library exposes raw access to the video memory.

http://www.ummon.eu/Linux/API/Devices/framebuffer.html

## Example

```go
package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/bfanger/framebuffer"
)

func main() {
	fb, err := framebuffer.Open("/dev/fb0")
	if err != nil {
		log.Fatalf("could not open framebuffer: %v", err)
	}
	defer fb.Close()
	info, err := fb.VarScreenInfo()
	if err != nil {
		log.Fatalf("could not read screen info: %v", err)
	}
	fmt.Printf(`
Fixed:
%+v

Variable:
%+v
`, fb.FixScreenInfo, info)
	rand.Read(fb.Buffer) // fill the buffer with noise
}
```

## Backstory

I got Raspberry Pi with a 2.8 inch TFT screen but could use opengl with resulted in only 15fps

To disable the blinking cursor:

```sh
echo 0 > /sys/class/graphics/fbcon/cursor_blink
```
