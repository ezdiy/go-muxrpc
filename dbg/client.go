/*
This file is part of go-muxrpc.

go-muxrpc is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

go-muxrpc is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with go-muxrpc.  If not, see <http://www.gnu.org/licenses/>.
*/

package muxrpc

import (
	"time"
	"os"
	"encoding/json"
	"net"
	"fmt"
	"testing"

	"github.com/ezdiy/go-muxrpc/codec"
	//"github.com/cryptix/go/logging/logtest"
	"github.com/cryptix/go/proc"
	"github.com/go-kit/kit/log"
)

func TestCall(t *testing.T) {
	//logger := log.NewLogfmtLogger(logtest.Logger("TestCall()", t))
	wr := log.NewSyncWriter(os.Stderr)
	logger := log.NewLogfmtLogger(wr)

	serv, err := proc.StartStdioProcess("node", wr, "client_test.js")
	//serv, err := proc.StartStdioProcess("node", wr, "client_test.js")
	if err != nil {
		t.Fatal(err)
	}

	//c := NewClient(logger, serv) //codec.Wrap(serv)) // debug.WrapRWC(serv)
	c := NewClient(logger, codec.Wrap(logger, serv)) // debug.WrapRWC(serv)
	var resp string
	fmt.Println("issuing call")
	err = c.Call("hello", &resp, "world", "bob")
	fmt.Println("run")
	if err != nil {
		t.Fatal(err)
	}
	if resp != "hello, world and bob!" {
		t.Fatal("wrong response:", resp)
	}
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)

}

func TestSource(t *testing.T) {
	//logger := log.NewLogfmtLogger(logtest.Logger("TestSyncSource()", t))
	//serv, err := proc.StartStdioProcess("node", logtest.Logger("client_test.js", t), "client_test.js")

	wr := log.NewSyncWriter(os.Stderr)
	logger := log.NewLogfmtLogger(wr)

	serv, err := proc.StartStdioProcess("node", wr, "client_test.js")

	if err != nil {
		t.Fatal(err)
	}
	c := NewClient(logger, codec.Wrap(logger, serv)) // debug.WrapRWC(serv)
	//c := NewClient(logger, serv) //codec.Wrap(logger,serv))
	resp := make(chan struct{ A int })

	go func() {
		c.Source("stuff", resp)
		close(resp)
	}()
	count := 0
	for range resp {
		//fmt.Printf("%#v\n", val)
		count++
	}
	if count != 4 {
		t.Fatal("Incorrect number of elements")
	}
	/*
		 // TODO: test values again
			sort.Ints(resp)
			for i := 0; i < 5; i++ {
				if resp[i] != i+1 {
					t.Errorf("resp missing: %d", resp[i])
				}
			}
	*/
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
}

func TestFullCall(t *testing.T) {
	p1, p2 := net.Pipe()
	//logger := log.NewLogfmtLogger(logtest.Logger("TestFull()", t))

	wr := log.NewSyncWriter(os.Stderr)
	logger := log.NewLogfmtLogger(wr)



	server := NewClient(logger, p1)
	client := NewClient(logger, p2)

	server.HandleCall("test", func(args json.RawMessage) interface{} {
		return "test"
	})
	go server.Handle()

	var resp string
	client.Call("test", &resp)

	if resp != "test" {
		t.Fatal("wrong response: ", resp)
	}
	time.Sleep(100 * time.Millisecond)
}

func TestFullSource(t *testing.T) {
	p1, p2 := net.Pipe()
	//logger := log.NewLogfmtLogger(logtest.Logger("TestFull()", t))
	wr := log.NewSyncWriter(os.Stderr)
	logger := log.NewLogfmtLogger(wr)



	server := NewClient(logger, p1)

	client := NewClient(logger, p2)

	server.HandleSource("test", func(args json.RawMessage) chan interface{} {
		stream := make(chan interface{}, 4)
		stream <- "a"
		stream <- "b"
		stream <- "c"
		stream <- "d"
		close(stream)
		return stream
	})
	go server.Handle()

	resp := make(chan string)
	go func() {
		err := client.Source("test", resp)
		if err != nil {
			t.Fatal(err)
		}
		close(resp)
	}()

	count := 0
	for range resp {
		//fmt.Printf("%#v\n", val)
		count++
	}
	if count != 4 {
		t.Fatal("Incorrect number of elements")
	}
	time.Sleep(100 * time.Millisecond)
}
