package simplekvs

type SimpleKVS struct {}

func NewSimpleKVS() *SimpleKVS {
  return &SimpleKVS{}
}

func (kvs *SimpleKVS) Set() {}

func (kvs *SimpleKVS) Get() {}

func (kvs *SimpleKVS) Update() {}

func (kvs *SimpleKVS) Delete() {}
