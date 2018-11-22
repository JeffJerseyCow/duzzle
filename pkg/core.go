package duzzle

import (
	"errors"
	"fmt"
	"github.com/jeffjerseycow/gdb-mi"
)

// DuzzleContext is a context structure that holds the internal state of
// the Duzzle library.
type DuzzleContext struct {
	pid            string
	arch           string
	gdb            *gdb.Gdb
	inChannel      chan map[string]interface{}
	onNotification gdb.NotificationCallback
}

// Global context variable to allow callback access to ctx.inChannel.
var ctx = DuzzleContext{}

// callback registers itself with Gdb.New(...) and is utlised when reading
// asynchronous messages from Gdb.
func callback(notification map[string]interface{}) {
	ctx.inChannel <- notification
}

// setPid searches through the messages returned by callback looking for the
// process startup notification. If found the PID is read and set within the
// context structure.
func (c *DuzzleContext) setPid() error {

	for {
		res := <-c.inChannel

		if res["type"] == "notify" && res["class"] == "thread-group-started" {
			if payload, err := getPayload(res, "payload"); err == nil {
				c.pid = payload["pid"].(string)
				break
			}
		}
	}

	return nil
}

// New configures the Duzzle library and returns a pointer to the initialised
// context struct.
func New(arch string) (*DuzzleContext, error) {
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
func (c *DuzzleContext) Exit() {
	c.gdb.Exit()
}

func (c *DuzzleContext) Connect(addr string, port string) error {
	res, err := c.gdb.Send(
		"target-select", "remote", fmt.Sprintf("%s:%s", addr, port))

	if err != nil {
		return errors.New(
			fmt.Sprintf("duzzle:Connect: Cannot connect to target %s:%s", addr, port))
	}

	if res["class"] == "connected" {
		c.setPid()
		return nil
	}

	return errors.New("duzzle:Connect: Unknown connection error")
}

func (c *DuzzleContext) Disconnect() error {

	if _, err := c.gdb.Send("target-detach"); err == nil {
		return nil
	}

	return errors.New("duzzle:Disconnect: Unknown disconnection error")
}

func (c *DuzzleContext) Breakpoint(addr uint64) error {
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

func (c *DuzzleContext) Continue() error {
	res, err := c.gdb.Send("exec-continue")

	if err != nil {
		return errors.New("duzzle:Continue: Cannot continue execution")
	}

	if res["class"].(string) == "running" {
		return nil
	}

	return errors.New("duzzle:Continue: Unknown continue error")
}

func (c *DuzzleContext) WaitBreak(addr uint64) error {

	for {
		res := <-c.inChannel
		if res["type"] == "exec" && res["class"] == "stopped" {
			if frame, err := getPayload(res, "payload", "frame"); err == nil && addrCmp(addr, frame["addr"].(string)) {
				break
			}
		}
	}

	return nil
}
