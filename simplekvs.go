package simplekvs

import "os"

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

func (kvs *SimpleKVS) Set() {}

func (kvs *SimpleKVS) Get() {}

func (kvs *SimpleKVS) Update() {}

func (kvs *SimpleKVS) Delete() {}
