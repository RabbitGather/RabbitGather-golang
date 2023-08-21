package lib

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/meowalien/RabbitGather-proto/go/share"
	"github.com/meowalien/go-meowalien-lib/errs"
)

func DecodeMessage[T any](encoding share.Encoding, data []byte, s *T) (err error) {
	switch encoding {
	case share.Encoding_GOB:
		err = gob.NewDecoder(bytes.NewReader(data)).Decode(s)
		if err != nil {
			err = errs.New(err)
			return
		}
	default:
		err = errs.New(ErrUnknownEncoding)
	}

	return
}

var ErrUnknownEncoding = errors.New("unknown encoding")
