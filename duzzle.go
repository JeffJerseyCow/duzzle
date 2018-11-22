package main

import (
	"fmt"
	"github.com/jeffjerseycow/duzzle/pkg"
)

var version = "undefined"

func main() {
	fmt.Println("Version", version)
	ctx, _ := duzzle.New("x86_64")

	ctx.Connect("127.0.0.1", "4444")

	fmt.Printf("Connected to %s:%s\n", "127.0.0.1", "4444")

	ctx.Breakpoint(0x555555559900)

	ctx.Continue()

	ctx.WaitBreak(0x555555559900)

	ctx.Disconnect()

	ctx.Exit()
}
