package simplekvs

import (
	"fmt"
	"os"
)

type Idx map[string]int

type SimpleKVS struct {
	f   *os.File
	idx Idx
}

func NewSimpleKVS(file string) (*SimpleKVS, error) {
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	idx := map[string]int{}
	return &SimpleKVS{
		f:   f,
		idx: idx,
	}, nil
}

type SimpleKVSError struct {
	msg string
	err error
}

func (e *SimpleKVSError) Error() string {
	return fmt.Sprintf("Error from SimpleKVS: %s (%s)", e.msg, e.err.Error())
}

func (e *SimpleKVSError) Unwrap() error {
	return e.err
}

func (kvs *SimpleKVS) Set(k string, v string) error {
	s, err := kvs.f.Stat()
	if err != nil {
		return &SimpleKVSError{err: err, msg: "Failed to get stat in Set"}
	}

	pos := s.Size()
	iniPos := pos

	// バリューの書き込み
	if n, err := kvs.f.WriteAt([]byte(v), pos); err == nil {
		pos += int64(n)
	} else {
		return &SimpleKVSError{err: err, msg: "Failed to write bytes in Set"}
	}

	// バリュー長の書き込み
	if n, err := kvs.f.WriteAt(
		[]byte([]uint8{uint8(len(v))}),
		pos,
	); err == nil {
		pos += int64(n)
	} else {
		return &SimpleKVSError{err: err, msg: "Failed to write bytes in Set"}
	}

	// キーの書き込み
	if n, err := kvs.f.WriteAt([]byte(k), pos); err == nil {
		pos += int64(n)
	} else {
		return &SimpleKVSError{err: err, msg: "Failed to write bytes in Set"}
	}

	// キー長の書き込み
	if n, err := kvs.f.WriteAt(
		[]byte([]uint8{uint8(len(k))}),
		pos,
	); err == nil {
		pos += int64(n)
	} else {
		return &SimpleKVSError{err: err, msg: "Failed to write bytes in Set"}
	}

	kvs.idx[k] = int(iniPos)

	return nil
}

func (kvs *SimpleKVS) Get(k string) (string, error) {
	return "", nil
}

func (kvs *SimpleKVS) Update(k string, v string) error {
	return nil
}

func (kvs *SimpleKVS) Delete(k string) error {
	return nil
}
