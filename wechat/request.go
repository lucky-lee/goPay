package wechat

import "github.com/lucky-lee/gutil/gStr"

//create payment request
func newPay() (res *ConfigPay) {
	res = &ConfigPay{}

	res.AppId = AppId()
	res.MchId = MchId()
	res.NonceStr = gStr.RandStr(32)

	return
}

//create app payment request
func NewPayApp() (res *ConfigPay) {
	res = newPay()

	res.setTypeApp()
	res.Source = SOURCE_APP

	return
}

//create applet payment request
func NewPayApplet() (res *ConfigPay) {
	res = newPay()

	res.setTypeJsapi()
	res.Source = SOURCE_APPLET

	return
}
