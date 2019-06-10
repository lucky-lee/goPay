package wechat

import (
	"fmt"
	"github.com/lucky-lee/gutil/gStr"
	"time"
)

//二次签名支付配置
type TwoSignApp struct {
	Appid      string `json:"appid" xml:"appid"`
	Noncestr   string `json:"noncestr" xml:"noncestr"`
	packAge    string `json:"package" xml:"package"`
	Partnerid  string `json:"partnerid" xml:"partnerid"`
	Prepayid   string `json:"prepayid" xml:"prepayid"`
	Timestamp  string `json:"timestamp" xml:"timestamp"`
	Sign       string `json:"sign" xml:"sign"`
	Packageval string `json:"packageval" xml:"packageval"`
}

//生成-二次签名-配置
func newTwoSignApp(prepayId string) TwoSignApp {
	var b TwoSignApp
	pVal := "Sign=WXPay"

	b.Appid = AppId()
	b.Partnerid = MchId()
	b.Prepayid = prepayId
	b.Noncestr = gStr.RandStr(32)
	b.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	b.packAge = pVal
	b.Packageval = pVal
	b.Sign = toSign(b)

	return b
}
