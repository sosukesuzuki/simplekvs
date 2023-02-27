package simplekvs_test

import (
	"testing"

	"github.com/sosukesuzuki/simplekvs"
)

func TestSimpleKVS(t *testing.T) {
	type Call struct {
		// "Set" | "Get" | "Delete"
		method  string
		args    []string
		wantRet string
	}

	tests := []struct {
		name  string
		calls []Call
	}{
		{
			name: "Setしていない値をGetすると空文字が返ってくる",
			calls: []Call{
				{
					method: "get",
					args: []string{"foo"},
					wantRet: "",
				},
			},
		},
		{
			name: "Setした値をGetできる",
			calls: []Call{
				{
					method:  "set",
					args:    []string{"foo", "bar"},
					wantRet: "",
				},
				{
					method:  "get",
					args:    []string{"foo"},
					wantRet: "bar",
				},
			},
		},
		{
			name: "同一のキーに対して複数回Setした場合、最新の値がGetできる",
			calls: []Call{
				{
					method:  "set",
					args:    []string{"foo", "bar1"},
					wantRet: "",
				},
				{
					method:  "set",
					args:    []string{"foo", "bar2"},
					wantRet: "",
				},
				{
					method:  "get",
					args:    []string{"foo"},
					wantRet: "bar2",
				},
			},
		},
		{
			name: "複数のキーをSetしたときも、Getできる",
			calls: []Call{
				{
					method:  "set",
					args:    []string{"foo1", "bar1"},
					wantRet: "",
				},
				{
					method:  "set",
					args:    []string{"foo2", "bar2"},
					wantRet: "",
				},
				{
					method:  "get",
					args:    []string{"foo1"},
					wantRet: "bar1",
				},
			},
		},
		{
			name: "複数のキーをSetしたときも、Getできる",
			calls: []Call{
				{
					method:  "set",
					args:    []string{"foo1", "bar1"},
					wantRet: "",
				},
				{
					method:  "set",
					args:    []string{"foo2", "bar2"},
					wantRet: "",
				},
				{
					method:  "get",
					args:    []string{"foo2"},
					wantRet: "bar2",
				},
			},
		},
		{
			name: "Deleteした値はGetできない",
			calls: []Call{
				{
					method:  "set",
					args:    []string{"foo", "bar"},
					wantRet: "",
				},
				{
					method:  "delete",
					args:    []string{"foo"},
					wantRet: "",
				},
				{
					method:  "get",
					args:    []string{"foo"},
					wantRet: "",
				},
			},
		},
	}

	for _, tt := range tests {
		kvs, err := simplekvs.NewSimpleKVS("dummy")
		if err != nil {
			t.Error("KVSインスタンスの作成に失敗")
		}
		defer kvs.Close()
		t.Run(tt.name, func(t *testing.T) {
			for _, call := range tt.calls {
				switch call.method {
				case "set":
					{
						err = kvs.Set(call.args[0], call.args[1])
						if err != nil {
							t.Errorf("k = %s と v = %s で呼び出したSetの呼び出しに失敗", call.args[0], call.args[1])
						}
						break
					}
				case "get":
					{
						ret, err := kvs.Get(call.args[0])
						if err != nil {
							t.Errorf("k = %s で呼び出した Get が失敗", call.args[0])
						}
						if call.wantRet != ret {
							t.Errorf("予期せぬ Set の返り値、期待する値: %s、実際の値: %s", call.wantRet, ret)
						}
						break
					}
				case "delete":
					{
						err = kvs.Delete(call.args[0])
						if err != nil {
							t.Errorf("k = %s で呼び出した Delete が失敗", call.args[0])
						}
						break
					}
				default:
					t.Errorf("メソッド %s はサポートされていません", call.method)
				}
			}
		})
	}
}
