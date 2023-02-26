package main

import (
	"fmt"
	"os"
	"strconv"
)

type Idx map[string]int

type SimpleKVS struct {
	f   *os.File
	idx Idx
}

func NewSimpleKVS(file string) (*SimpleKVS, error) {
	f, err := os.Create(file)
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

func Close(kvs *SimpleKVS) error {
	err := kvs.f.Close()
	return err
}

func (kvs *SimpleKVS) Set(k string, v string) error {
	s, err := kvs.f.Stat()
	if err != nil {
		return &SimpleKVSError{err: err, msg: "Failed to get stat in Set"}
	}

	pos := s.Size()
	iniPos := pos

	// バリュー長の書き込み
	if n, err := kvs.f.WriteAt(
		[]byte(fmt.Sprintf("%d", len(v))),
		pos,
	); err == nil {
		pos += int64(n)
	} else {
		return &SimpleKVSError{err: err, msg: "Failed to write bytes in Set"}
	}

	// バリューの書き込み
	if n, err := kvs.f.WriteAt([]byte(v), pos); err == nil {
		pos += int64(n)
	} else {
		return &SimpleKVSError{err: err, msg: "Failed to write bytes in Set"}
	}

	// キー長の書き込み
	if n, err := kvs.f.WriteAt(
		[]byte(fmt.Sprintf("%d", len(k))),
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

	kvs.idx[k] = int(iniPos)

	return nil
}

func (kvs *SimpleKVS) Get(k string) (string, error) {
	pos, ok := kvs.idx[k]

	// インデックスの中にキーに対応するポジションがなければエラー
	if !ok {
		return "", &SimpleKVSError{
			err: nil,
			msg: fmt.Sprintf("Failed to find key %s in Get", k),
		}
	}

	// ファイルをposまでSeekする
	s, err := kvs.f.Seek(int64(pos), 0)
	if err != nil {
		return "", &SimpleKVSError{
			err: err,
			msg: "Failed to seek file in Get",
		}
	}
	fmt.Printf("Seek %d\n", s)

	// 後続のバリューの長さをvalue_lengthに読み込む
	// これがバリューの長さに該当する
	value_length := make([]byte, 1)
	n1, err := kvs.f.Read(value_length)
	if err != nil {
		return "", &SimpleKVSError{
			err: err,
			msg: "Failed to read value_length in Get",
		}
	}
	fmt.Printf("Read %d\n", n1)

	// value_length分をvalueに読み込む
	value_length_i, err := strconv.Atoi(string(value_length))
	value := make([]byte, value_length_i)
	n2, err := kvs.f.Read(value)
	if err != nil {
		return "", &SimpleKVSError{
			err: err,
			msg: "Failed to read value in Get",
		}
	}
	fmt.Printf("Read %d\n", n2)

	return string(value), nil
}

func (kvs *SimpleKVS) Delete(k string) error {
	return nil
}
