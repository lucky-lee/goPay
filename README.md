# golang版 支付工具类 目前只有微信支付,后期会接入支付宝等其他支付

## 本项目采用go mod 需要 golang v1.11.+

##安装
```go
go get -u github.com/lucky-lee/goPay
```

## 微信 app 统一下单
```go
func PayWxApp() {
	//设置基本信息
	wechat.SetAppId("appId123456")
	wechat.SetMchId("mchId123456")
	wechat.SetApiKey("apiKey123456")

	//生成预支付订单
	res := wechat.NewPay().
		SetTypeApp().
		SetBody("test").
		SetClientIp("127.0.0.1").
		SetNotifyUrl("http://www.xxxxx.com").
		SetOutTradeNo("123456789").
		SetTotalFree(0.01).CreatePayOrder()

	fmt.Println(gJson.Encode(res))
}
```
## 微信 小程序 统一下单
```go
func PayWxApplet() {
	//设置基本信息
	wechat.SetAppId("appId123456")
	wechat.SetMchId("mchId123456")
	wechat.SetApiKey("apiKey123456")

	//生成预支付订单
	res := wechat.NewPay().
		SetTypeJsapi().
		SetBody("test").
		SetClientIp("127.0.0.1").
		SetNotifyUrl("http://www.xxxxx.com").
		SetOutTradeNo("123456789").
		SetTotalFree(0.01).CreatePayOrder()

	fmt.Println(gJson.Encode(res))
}
```