package wechat

import (
	"encoding/json"
	"testing"
)

func TestSignJsapi(t *testing.T) {
	res := newTwoSignApplet("test")
	bs, _ := json.Marshal(res)

	t.Log(string(bs))
}
