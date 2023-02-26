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

	// バリュー長の書き込み
	if n, err := kvs.f.WriteAt(
		[]byte([]uint8{uint8(len(v))}),
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
		[]byte([]uint8{uint8(len(k))}),
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
	pos := kvs.idx[k]
	// インデックスの中にキーに対応するポジションがなければエラー
	if pos == 0 {
		return "", &SimpleKVSError{
			err: nil,
			msg: fmt.Sprintf("Failed to find key %s in Get", k),
		}
	}

	// ファイルをposまでSeekする
	s1, err := kvs.f.Seek(int64(pos), 0)
	_ = s1
	if err != nil {
		return "", &SimpleKVSError{
			err: err,
			msg: "Failed to seek file in Get",
		}
	}

	// 1バイト分をvalue_lengthに読み込む
	// これがバリューの長さに該当する
	value_length := make([]byte, 1)
	n1, err := kvs.f.Read(value_length)
	if err != nil {
		return "", &SimpleKVSError{
			err: err,
			msg: "Failed to read file in Get",
		}
	}

	// value_lengthのバイト数(1バイト)分だけSeekする
	s2, err := kvs.f.Seek(int64(n1), 2)
	_ = s2
	if err != nil {
		return "", &SimpleKVSError{
			err: err,
			msg: "Failed to seek file in Get",
		}
	}

	// value_length分をvalueに読み込む
	value := make([]byte, value_length[0])
	n2, err := kvs.f.Read(value)
	_ = n2
	if err != nil {
		return "", &SimpleKVSError{
			err: err,
			msg: "Failed to read file in Get",
		}
	}

	return string(value), nil
}

func (kvs *SimpleKVS) Update(k string, v string) error {
	return nil
}

func (kvs *SimpleKVS) Delete(k string) error {
	return nil
}
