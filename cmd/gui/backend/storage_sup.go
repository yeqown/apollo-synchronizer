package backend

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"os"

	"github.com/pkg/errors"
)

type ext string

const (
	_ext_binary ext = "bin"
	_ext_json   ext = "json"
)

func read(fp string, rcvr interface{}, _ext ext) (err error) {
	fd, err2 := os.Open(fp)
	if err2 != nil {
		return errors.Wrap(err2, "open file failed")
	}
	defer fd.Close()

	switch _ext {
	case _ext_binary:
		err = binary.Read(fd, binary.BigEndian, rcvr)
	case _ext_json:
		var bytes []byte
		if bytes, err = io.ReadAll(fd); err != nil {
			return errors.Wrap(err, "read file failed")
		}

		err = json.Unmarshal(bytes, rcvr)
		err = errors.Wrap(err, "unmarshal json failed")
	}

	if err != nil {
		return err
	}

	return nil
}

// save file overwrite.
func save(fp string, data interface{}, _ext ext) error {
	fd, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	switch _ext {
	case _ext_binary:
		binary.Write(fd, binary.BigEndian, data)
	case _ext_json:
		bytes, err := json.Marshal(data)
		if err != nil {
			return errors.Wrap(err, "marshal json failed")
		}
		_, err = fd.Write(bytes)
	}

	if err != nil {
		return err
	}

	return nil
}
