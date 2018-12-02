package duzzle

import (
	"errors"
	"fmt"
	"github.com/JeffJerseyCow/gdb-mi"
)

// New configures the Duzzle library and returns a pointer to the initialised
// context struct.
func New(arch string) (*duzzleContext, error) {
	ctx.arch = arch
	ctx.inChannel = make(chan map[string]interface{}, 1024)
	ctx.onNotification = callback
	gdb, err := gdb.New(callback)

	if err != nil {
		return nil, errors.New("duzzle:New: Cannot initialise GDB/MI")
	}

	ctx.gdb = gdb
	return &ctx, nil
}

// Exit tears down the Gdb connection. Call Disconnect first to gracefully exit
// from the target process.
func (c *duzzleContext) Exit() {
	c.gdb.Exit()
}

func (c *duzzleContext) Connect(addr string, port string) error {
	res, err := c.gdb.Send(
		"target-select", "remote", fmt.Sprintf("%s:%s", addr, port))

	if err != nil {
		return errors.New(
			fmt.Sprintf("duzzle:Connect: Cannot connect to target %s:%s", addr, port))
	}

	if res["class"] == "connected" {
		setPid(c)
		return nil
	}

	return errors.New("duzzle:Connect: Unknown connection error")
}

func (c *duzzleContext) Disconnect() error {

	if _, err := c.gdb.Send("target-detach"); err == nil {
		return nil
	}

	return errors.New("duzzle:Disconnect: Unknown disconnection error")
}

func (c *duzzleContext) Breakpoint(addr uint64) error {
	res, err := c.gdb.Send("break-insert", fmt.Sprintf("*0x%x", addr))

	if err != nil {
		return errors.New("duzzle:Breakpoint: Cannot insert breakpoint")
	}

	if bkpt, err := getPayload(res, "payload", "bkpt"); err == nil {
		if addrCmp(addr, bkpt["addr"].(string)) {
			return nil
		}
	}

	return errors.New("duzzle:Breakpoint: Unknown breakpoint error")
}

func (c *duzzleContext) Continue() error {
	res, err := c.gdb.Send("exec-continue")

	if err != nil {
		return errors.New("duzzle:Continue: Cannot continue execution")
	}

	if res["class"].(string) == "running" {
		return nil
	}

	return errors.New("duzzle:Continue: Unknown continue error")
}

// WaitBreak loops reading debugging information from gdbserver until execution
// stop on the breakpoint specified by addr. Ensure a breakpoint has been set
// using Breakpoint prior to calling this function.
func (c *duzzleContext) WaitBreak(addr uint64) error {

	for {
		res := <-c.inChannel
		if res["type"] == "exec" && res["class"] == "stopped" {
			if frame, err := getPayload(res, "payload", "frame"); err == nil {
				if addrCmp(addr, frame["addr"].(string)) {
					break
				}
			}
		}
	}

	return nil
}

// Map downloads the remote process maps file and extracts segment data with
// each being stored in a map type.
func (c *duzzleContext) Map() ([]map[string]interface{}, error) {
	srcFile := fmt.Sprintf("/proc/%s/maps", c.pid)
	mapFile, _ := downloadFile(c, srcFile, "maps")
	return parseMap(mapFile)
}
