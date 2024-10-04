package packetti

import (
	"bytes"
	"fmt"
	"io"
)

const (
	slipEnd    = 0xC0
	slipEsc    = 0xDB
	slipEscEnd = 0xDC
	slipEscEsc = 0xDD
)

func readSlipPacket(r io.Reader) ([]byte, error) {
	buf := &bytes.Buffer{}
	readBuf := make([]byte, 1)

	for {
		n, err := r.Read(readBuf)
		if n == 0 {
			continue
		}

		if err != nil {
			return nil, fmt.Errorf("read next char: %w", err)
		}

		switch readBuf[0] {
		case slipEnd:
			if buf.Len() == 0 {
				continue
			}

			return buf.Bytes(), nil

		case slipEsc:
			n, err = r.Read(readBuf)

			if n == 0 {
				continue
			}

			if err != nil {
				return nil, fmt.Errorf("read next char: %w", err)
			}

			/* if "c" is not one of these two, then we
			 * have a protocol violation.  The best bet
			 * seems to be to leave the byte alone and
			 * just stuff it into the packet
			 */
			switch readBuf[0] {
			case slipEscEnd:
				readBuf[0] = slipEnd
			case slipEscEsc:
				readBuf[0] = slipEsc
			}
		}

		buf.WriteByte(readBuf[0])
	}
}
