package wechat

import (
	"encoding/json"
	"testing"
)

func TestSignJsapi(t *testing.T) {
	res := newTwoSignJsapi("test")
	bs, _ := json.Marshal(res)

	t.Log(string(bs))
}
