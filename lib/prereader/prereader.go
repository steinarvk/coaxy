package prereader

import (
	"io"
)

type reader struct {
	remainingBytes []byte
	underlying     io.Reader
	eof            bool
}

func (r *reader) Read(data []byte) (int, error) {
	if r.eof {
		return 0, io.EOF
	}

	var nread int

	if len(r.remainingBytes) > 0 {
		n := copy(data, r.remainingBytes)
		r.remainingBytes = r.remainingBytes[n:]
		nread += n
	}

	if nread < len(data) && r.underlying != nil {
		n, err := r.underlying.Read(data[nread:])
		nread += n
		if err != nil && err != io.EOF {
			return nread, err
		}
		if err == io.EOF {
			r.eof = true
		}
	}

	if nread == 0 && r.underlying == nil {
		return 0, io.EOF
	}

	return nread, nil
}

func Preread(r io.Reader, limit int) ([]byte, io.Reader, error) {
	if limit <= 0 {
		return nil, r, nil
	}

	buf := make([]byte, limit)
	eof := false
	var nread int

	for {
		n, err := r.Read(buf[nread:])

		nread += n

		if err == io.EOF {
			eof = true
			break
		}
		if err != nil {
			return nil, nil, err
		}

		if nread >= limit {
			break
		}
	}

	data := buf[:nread]
	rv := &reader{
		remainingBytes: data,
	}
	if !eof {
		rv.underlying = r
	}

	return data, rv, nil
}
