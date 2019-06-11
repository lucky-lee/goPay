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
	Sign       string `json:"sign" xml:"-"`
	Packageval string `json:"packageval" xml:"-"`
}

//生成-二次签名-配置
func newTwoSignApp(prepayId string) (res TwoSignApp) {
	pVal := "Sign=WXPay"

	res.Appid = AppId()
	res.Partnerid = MchId()
	res.Prepayid = prepayId
	res.Noncestr = gStr.RandStr(32)
	res.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	res.packAge = pVal
	res.Packageval = pVal
	res.Sign = toSign(res)

	return
}
