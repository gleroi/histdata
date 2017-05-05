package histdata

import (
	"encoding/binary"
	"io"
	"strings"
	"time"
)

type Entry struct {
	Attribute string
	Date      time.Time
	Value     float64
}

type Reader struct {
	input     io.Reader
	attribute string
}

func NewReader(input io.Reader) (*Reader, error) {
	r := Reader{
		input: input,
	}
	err := r.init()
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (r *Reader) init() error {
	var length int16
	err := binary.Read(r.input, binary.BigEndian, &length)
	if err != nil {
		return err
	}
	desc := make([]byte, length)
	_, err = r.input.Read(desc)
	if err != nil {
		return err
	}
	r.attribute = strings.Replace(string(desc), "attribute://", "", -1)
	return nil
}

func (r *Reader) Name() string {
	return r.attribute
}

func (r *Reader) Read() (Entry, error) {
	t, err := readTime(r)
	if err != nil {
		return Entry{}, err
	}
	v, err := readDouble(r)
	if err != nil {
		return Entry{}, err
	}
	return Entry{
		Attribute: r.attribute,
		Date:      t,
		Value:     v,
	}, nil
}

func readTime(r *Reader) (time.Time, error) {
	var val int64
	err := binary.Read(r.input, binary.BigEndian, &val)
	if err != nil {
		return time.Time{}, err
	}
	const nanoseconds = 1000000000
	t := time.Unix(val/nanoseconds, val%nanoseconds)
	return t, nil
}

func readDouble(r *Reader) (float64, error) {
	var val float64
	err := binary.Read(r.input, binary.BigEndian, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}
