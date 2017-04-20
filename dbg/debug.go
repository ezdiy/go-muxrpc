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

package codec

import (
	//"time"
	"io"
	"fmt"

	"github.com/go-kit/kit/log"
)

// Wrap decodes every packet that passes through it and logs it
func Wrap(l log.Logger, rwc io.ReadWriteCloser) io.ReadWriteCloser {
	prout, pwout := io.Pipe()
	lin := log.With(l, "unit", "Recvd")
	go func() {
		w := NewWriter(pwout)
		r := NewReader(rwc)
		for {
			pkt, err := r.ReadPacket()
			//time.Sleep(1000 * time.Millisecond)
			if err != nil {
				if err != io.EOF {
					lin.Log("action", "ReadPacket", "error", err)
					prout.CloseWithError(err)
				}
				return
			}
			if err := w.WritePacket(pkt); err != nil {
				fmt.Println("write failed?")
				lin.Log("action", "WritePacket", "error", err)
				prout.CloseWithError(err)
				return
			}
			lin.Log("pkt", pkt)
		}
	}()

	prin, pwin := io.Pipe()
	w := NewWriter(rwc)
	lin2 := log.With(l, "unit", "Sent")
	go func() {
		r := NewReader(prin)
		for {
			pkt, err := r.ReadPacket()
			//time.Sleep(1000 * time.Millisecond)
			if err != nil {
				if err != io.EOF {
					lin2.Log("action", "ReadPacket", "error", err)
					prin.CloseWithError(err)
				}
				return
			}
			if err := w.WritePacket(pkt); err != nil {
				lin2.Log("action", "WritePacket", "error", err)
				prin.CloseWithError(err)
				return
			}
			lin2.Log("pkt", pkt)
		}
	}()
	return struct {
		io.Reader
		io.Writer
		io.Closer
	}{Reader: prout, Writer: pwin, Closer: w}
}
