# SimpleKVS

[データ指向アプリケーションデザイン](https://www.oreilly.co.jp/books/9784873118703/)の 3 章を読んで学習のために作ってみたシンプルな KVS です。

## インデックス

ただのインメモリマップを保持しています。キーから、そのキーに対応するバリューが存在するファイル内のオフセットへのマップです。

## ログフォーマット

次のようなキーとバリューがあるとします。

| key    | value  |
| ------ | ------ |
| foo    | bar    |
| foofoo | barbar |

上から順番に書き込んでいった場合、次のように文字列ファイルに書き込まれています。

```
3foo3bar6foofoo6barbar
```

## API

### SimpleKVS

```go
type SimpleKVS struct {}
```

### NewSimpleKVS

```go
func NewSimpleKVS(file string) (*SimpleKVS, error)
```

### func Close(kvs \*SimpleKVS)

```go
func (kvs *SimpleKVS) Close() error
```

### func (kvs \*SimpleKVS) Set(k string, v string)

```go
func (kvs *SimpleKVS) Set(k string, v string) error
```

### func (kvs \*SimpleKVS) Get(k string)

```go
func (kvs *SimpleKVS) Get(k string)
```

### func (kvs \*SimpleKVS) Delete(k string) error

```go
func (kvs *SimpleKVS) Delete(k string) error
```
