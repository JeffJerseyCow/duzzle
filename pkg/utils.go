package duzzle

import (
  "fmt"
  "errors"
  "strconv"
)

// getPayload take a variadic string and iterates through nested payloads
// searching for the matching strings specified by strPayloads -- if found
// the requested payload is returned. Use the following to parse a "bkpt"
// payload within a parent "payload" in the message res.
//  GetPayload(res, "payload", "bkpt")
func getPayload(res map[string]interface{}, strPayloads ...string) (
  map[string]interface{}, error) {

  for _, strPayload := range strPayloads {
    if payload, ok := res[strPayload]; ok {
      res = payload.(map[string]interface{})
    } else {
      return nil, errors.New(fmt.Sprintf(
        "duzzle:GetPayload: Unknown payload name '%s'", strPayload))
    }
  }

  return res, nil
}

// addrCmp compares a unsigned integer value with a hexidecimal string and
// returns true if they match or false if they don't.
func addrCmp(addr uint64, strAddr string) (bool) {

  if strAddr, err := strconv.ParseUint(strAddr, 0, 64);
     err == nil && addr == strAddr {
    return true
  } else {
    return false
  }
}
