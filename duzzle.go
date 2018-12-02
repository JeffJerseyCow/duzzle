package main

import (
	"fmt"
	"github.com/JeffJerseyCow/duzzle/core"
)

var version = "undefined"

func main() {
	IPV4Addr := "127.0.0.1"
	Port := "4444"
	var Breakpoint uint64 = 0x555555559900
	ctx, _ := duzzle.New("x86_64")
	defer ctx.Exit()
	ctx.Connect(IPV4Addr, Port)
	defer ctx.Disconnect()
	fmt.Println(fmt.Sprintf("Connected to: %s:%s", IPV4Addr, Port))
	ctx.Breakpoint(Breakpoint)
	fmt.Println(fmt.Sprintf("Set breakpoing: 0x%x", Breakpoint))
	ctx.Continue()
	fmt.Println("Continuting execution")
	ctx.WaitBreak(Breakpoint)
	fmt.Println(fmt.Sprintf("Hit breakpoint: 0x%x", Breakpoint))
	segmentMap, _ := ctx.Map()

	for _, segment := range segmentMap {
		fmt.Println(segment)
	}

}
