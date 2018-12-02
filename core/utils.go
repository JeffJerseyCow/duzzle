package duzzle

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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
func addrCmp(addr uint64, strAddr string) bool {

	if strAddr, err := strconv.ParseUint(strAddr, 0, 64); err == nil && addr == strAddr {
		return true
	} else {
		return false
	}
}

// parsePerms takes a string containing the permissions of the memory segment
// and returns a UNIX permission byte
func parsePerms(permStr string) byte {
	var perms byte = 0

	if strings.Contains(permStr, "r") {
		perms |= 0x4
	}

	if strings.Contains(permStr, "w") {
		perms |= 0x2
	}

	if strings.Contains(permStr, "x") {
		perms |= 0x1
	}

	return perms
}

// parseSegment takes the results of a string slice containing the relevant
// memory segment and parses it into a map[string]interface{}
func parseSegment(segment []string) map[string]interface{} {
	start := fmt.Sprintf("0x%s", segment[1])
	end := fmt.Sprintf("0x%s", segment[2])
	startUint, _ := strconv.ParseUint(start, 0, 64)
	endUint, _ := strconv.ParseUint(end, 0, 64)
	size := endUint - startUint
	perms := parsePerms(segment[3])
	return map[string]interface{}{
		"start": start,
		"end":   end,
		"size":  size,
		"perms": perms,
	}
}

// parseMap takes a filepath string to the relevant Linux maps file, reads the
// contents and returns a slice of segments
func parseMap(mapFile string) ([]map[string]interface{}, error) {
	var segmentMap []map[string]interface{}
	file, err := os.Open(mapFile)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("duzzle:parseMap: Cannot open '%s'", mapFile))
	}
	defer file.Close()

	r, _ := regexp.Compile("^([0-9a-fA-F]*)-([0-9a-fA-F]*)\\s*([rwxaps-]*)\\s")
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		segment := r.FindStringSubmatch(scanner.Text())

		if len(segment) != 4 {
			return nil, errors.New("duzzle:parseMap: Incorrect length when parsing map")
		}

		segmentMap = append(segmentMap, parseSegment(segment))
	}

	return segmentMap, nil
}
