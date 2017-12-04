package sse

import (
	"io"
)

func paddedCopy(dst io.Writer, src io.Reader, prefix, suffix string) (int64, error) {
	return paddedCopyBuffer(dst, src, prefix, suffix, make([]byte, 32*1024))
}

func paddedCopyBuffer(dst io.Writer, src io.Reader, prefix, suffix string, buf []byte) (written int64, err error) {
	paddingLen := len(prefix) + len(suffix)
	startOffset := len(prefix)
	stopOffset := len(buf) - len(suffix)

	copy(buf, prefix)

	for {
		nr, er := src.Read(buf[startOffset:stopOffset])
		if nr > 0 {
			copy(buf[startOffset+nr:], suffix)

			nw, ew := dst.Write(buf[:nr+paddingLen])
			if nw > 0 {
				written += int64(nw)
			}

			if ew != nil {
				err = ew
				break
			}

			if (nr + paddingLen) != nw {
				err = io.ErrShortWrite
				break
			}
		}

		if er != nil {
			if er != io.EOF {
				err = er
			}

			break
		}
	}

	return
}
