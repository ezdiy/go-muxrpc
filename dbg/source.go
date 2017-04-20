package main
import (
	"fmt"
	"os"
	"encoding/json"
	"github.com/ezdiy/go-muxrpc"
	"github.com/cryptix/go/proc"
	"github.com/go-kit/kit/log"
)
func main() {
	wr := log.NewSyncWriter(os.Stderr)
	logger := log.NewLogfmtLogger(wr)

	peer, _ := proc.StartStdioProcess("node", wr, "source.js")
	mux := muxrpc.NewClient(logger, peer)
	mux.HandleSource("stuff", func(rm json.RawMessage) chan interface{} {
		c := make(chan interface{})
		fmt.Println("params are ",string(rm[:]))
		go func() {
			for i:=0; i < 5; i++ {
				c  <- i
			}
		}()
		return c
	})
	for {
		mux.Handle()
	}
}
