package crawler

import (
	"bytes"
	"encoding/binary"
	"compress/gzip"
	"fmt"
	"io/ioutil"
)

func ParseGzip(data []byte, handleErr bool) ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		//with error
		fmt.Printf("[ParseGzip %t %d] NewReader error: %v, maybe data is ungzip: [%s]\n", handleErr, len(data), err)
		if handleErr {
			errHandler(data)
		}
		return nil, err
	} else {
		defer r.Close()
		undatas, err := ioutil.ReadAll(r)
		if err != nil {
			//with error
			fmt.Printf("[ParseGzip %t %d] ioutil.ReadAll error: %v: [%s]\n", handleErr, len(data), err)
			if handleErr {
				errHandler(data)
			}
			return nil, err
		} else {
			//buffer.Reset()
			return undatas, nil
		}
	}
}
var buffer bytes.Buffer
func errHandler(data []byte) {
	buffer.Write(data)
	msg, err := ParseGzip(buffer.Bytes(), false)
	if err == nil {
		fmt.Println("!!!!!!", string(msg[:]))
	}
}
