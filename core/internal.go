package duzzle

import (
	"errors"
	"fmt"
	"github.com/JeffJerseyCow/gdb-mi"
	"io/ioutil"
	"path"
)

// duzzleContext is a context structure that holds the internal state of
// the Duzzle library.
type duzzleContext struct {
	pid            string
	arch           string
	gdb            *gdb.Gdb
	inChannel      chan map[string]interface{}
	onNotification gdb.NotificationCallback
}

// Global context variable to allow callback access to ctx.inChannel.
var ctx = duzzleContext{}

// callback registers itself with Gdb.New(...) and is utlised when reading
// asynchronous messages from Gdb.
func callback(notification map[string]interface{}) {
	ctx.inChannel <- notification
}

// setPid searches through the messages returned by callback looking for the
// process startup notification. If found the PID is read and set within the
// context structure.
func setPid(c *duzzleContext) error {

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

// downloadFile takes a file path located on the remote system and transfers
// it to the local system within a temporary directory. dstFileName is its
// name within the context of the temporary directory.
func downloadFile(c *duzzleContext, srcFilePath, dstFileName string) (string, error) {
	tmpDir, _ := ioutil.TempDir("", fmt.Sprintf("duzzle%s_", c.pid))
	dstFilePath := path.Join(tmpDir, dstFileName)
	res, err := c.gdb.Send("target-file-get", srcFilePath, dstFilePath)

	if err != nil {
		return "", errors.New("duzzle:downloadFile: Cannot download file")
	}

	if res["class"] == "done" {
		return dstFilePath, nil
	}

	return "", errors.New("duzzle:downloadFile: Unknown file download error")
}
