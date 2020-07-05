package example

import (
	"fmt"
	"github.com/lucky-lee/goPay/wechat"
	"github.com/lucky-lee/gutil/gJson"
)

//微信 app 统一下单
func PayWxApp() {
	//设置基本信息
	wechat.SetAppId("appId123456")
	wechat.SetMchId("mchId123456")
	wechat.SetApiKey("apiKey123456")

	//生成预支付订单
	res := wechat.NewPayApp().
		SetBody("test").
		SetClientIp("127.0.0.1").
		SetNotifyUrl("http://www.xxxxx.com").
		SetOutTradeNo("123456789").
		SetTotalFree(0.01).CreatePayOrder()

	fmt.Println(gJson.Encode(res))
}

//微信 小程序 统一下单
func PayWxApplet() {
	//设置基本信息
	wechat.SetAppId("appId123456")
	wechat.SetMchId("mchId123456")
	wechat.SetApiKey("apiKey123456")

	//生成预支付订单
	res := wechat.NewPayApplet().
		SetBody("test").
		SetClientIp("127.0.0.1").
		SetNotifyUrl("http://www.xxxxx.com").
		SetOutTradeNo("123456789").
		SetTotalFree(0.01).CreatePayOrder()

	fmt.Println(gJson.Encode(res))
}
