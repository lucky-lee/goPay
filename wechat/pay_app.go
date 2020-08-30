package wechat

import (
	"fmt"
	"github.com/lucky-lee/gutil/gStr"
	"time"
)

//二次签名支付配置
type TwoSignApp struct {
	AppId      string `json:"app_id" xml:"appid"`
	NonceStr   string `json:"nonce_str" xml:"noncestr"`
	packAge    string `json:"package" xml:"package"`
	PartnerId  string `json:"partner_id" xml:"partnerid"`
	PrepayId   string `json:"prepay_id" xml:"prepayid"`
	Timestamp  string `json:"timestamp" xml:"timestamp"`
	Sign       string `json:"sign" xml:"-"`
	PackageVal string `json:"package_val" xml:"-"`
}

//生成-二次签名-配置
func newTwoSignApp(orderId, prepayId string) (res TwoSignApp) {
	pVal := "Sign=WXPay"

	res.AppId = AppId()
	res.PartnerId = MchId()
	res.PrepayId = prepayId
	res.NonceStr = gStr.RandStr(32)
	res.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	res.packAge = pVal
	res.PackageVal = pVal
	res.Sign = toSign(res)

	return
}
