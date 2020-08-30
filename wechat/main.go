package wechat

import (
	"encoding/xml"
	"fmt"
	"github.com/lucky-lee/gutil/gCurl"
	"github.com/lucky-lee/gutil/gLog"
)

//支付配置 -按照ascii 排序
type ConfigPay struct {
	AppId      string `xml:"appid"`                //必传-应用id
	Attach     string `xml:"attach"`               //非必传-附加数据
	Body       string `xml:"body"`                 //必传-商品描述
	Detail     string `xml:"detail"`               //非必传-商品详情
	DeviceInfo string `xml:"device_info"`          //非必传-设备号
	FeeType    string `xml:"fee_type"`             //非必传-标价币种 符合ISO 4217标准的三位字母代码，默认人民币：CNY
	GoodsTag   string `xml:"goods_tag"`            //非必传-订单优惠标记，使用代金券或立减优惠功能时需要的参数
	LimitPay   string `xml:"limit_pay"`            //非必传-上传此参数no_credit--可限制用户不能使用信用卡支付
	MchId      string `xml:"mch_id"`               //必传-商户号
	NonceStr   string `xml:"nonce_str"`            //必传-随机字符串
	NotifyUrl  string `xml:"notify_url"`           //必传-通知地址
	Openid     string `xml:"openid" json:"openid"` //非必传-trade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识
	OutTradeNo string `xml:"out_trade_no"`         //必传-商户订单号

	//非必传-trade_type=NATIVE时，此参数必传。此参数为二维码中包含的商品ID，商户自行定义。
	ProductId string `xml:"product_id"`

	//非必传-Y，传入Y时，支付成功消息和支付详情页将出现开票入口。需要在微信支付商户平台或微信公众平台开通电子发票功能，传此字段才可生效
	Receipt string `xml:"receipt"`

	//该字段常用于线下活动时的场景信息上报，支持上报实际门店信息，商户也可以按需求自己上报相关信息。该字段为JSON对象数据，
	//对象格式为{"store_info":{"id": "门店ID","name": "名称","area_code": "编码","address": "地址" }}
	//SceneInfo  string `xml:"scene_info"`

	Sign     string `xml:"sign"`             //必传-签名
	SignType string `xml:"sign_type"`        //非必传-签名类型，默认为MD5，支持HMAC-SHA256和MD5。
	ClientIp string `xml:"spbill_create_ip"` //必传-终端IP

	//非必传-订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010
	TimeExpire string `xml:"time_expire"`

	//非必传-订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010
	TimeStart string `xml:"time_start"`
	TotalFee  string `xml:"total_fee"`  //必传-总金额
	TradeType string `xml:"trade_type"` //必传-交易类型

	//自定义字段
	Source uint8 `xml:"-"` //来源
}

//设置-商品描述
func (self *ConfigPay) SetBody(val string) *ConfigPay {
	self.Body = val

	return self
}

//设置-商户订单号
func (self *ConfigPay) SetOutTradeNo(val string) *ConfigPay {
	self.OutTradeNo = val

	return self
}

//设置-总金额
func (self *ConfigPay) SetTotalFree(val float64) *ConfigPay {
	self.TotalFee = fmt.Sprintf("%.0f", val*100)

	return self
}

//设置-通知地址
func (self *ConfigPay) SetNotifyUrl(val string) *ConfigPay {
	self.NotifyUrl = val

	return self
}

//设置-终端ip
func (self *ConfigPay) SetClientIp(val string) *ConfigPay {
	self.ClientIp = val

	return self
}

//设置-openid
func (self *ConfigPay) SetOpenid(val string) *ConfigPay {
	self.Openid = val

	return self
}

//设置 交易类型 app
func (self *ConfigPay) setTypeApp() *ConfigPay {
	self.TradeType = TYPE_APP

	return self
}

//设置 交易类型 jsapi
func (self *ConfigPay) setTypeJsapi() *ConfigPay {
	self.TradeType = TYPE_JSAPI

	return self
}

//设置 交易类型 h5
func (self *ConfigPay) setTypeH5() *ConfigPay {
	self.TradeType = TYPE_H5

	return self
}

//设置 交易类型 native
func (self *ConfigPay) setTypeNative() *ConfigPay {
	self.TradeType = TYPE_NATIVE

	return self
}

//生成-预付订单
func (s ConfigPay) CreatePayOrder() (res interface{}) {
	res = ResEmpty{}
	resPay, err := s.requestWx()
	if err != nil { //解析错误
		return ResEmpty{}
	}
	//判断是否成功调用微信预支付接口
	if resPay.ReturnCode == CODE_SUCCESS && resPay.ResultCode == CODE_SUCCESS {
		switch s.Source {
		case SOURCE_APP:
			res = newTwoSignApp(s.OutTradeNo, resPay.PrepayId)
		case SOURCE_APPLET:
			res = newTwoSignApplet(s.OutTradeNo, resPay.PrepayId)
			//case jsapi支付: //TODO 需要实现

			//case h5支付: //TODO 需要实现

			//case native支付: //TODO 需要实现

		}

	}
	return
}

//app 支付订单
func (s ConfigPay) PayAppOrder() (res TwoSignApp) {
	resPay, err := s.requestWx()
	if err != nil {
		gLog.E("解析json错误", err.Error())
		return
	}
	res = newTwoSignApp(s.OutTradeNo, resPay.PrepayId)
	return
}

//小程序 支付订单
func (s ConfigPay) PayAppletOrder() (res TwoSignApplet) {
	resPay, err := s.requestWx()
	if err != nil {
		gLog.E("解析json错误", err.Error())
		return
	}
	res = newTwoSignApplet(s.OutTradeNo, resPay.PrepayId)
	return
}

//请求微信
func (s ConfigPay) requestWx() (res ResPay, err error) {
	s.Sign = toSign(s)     //生成签名字符串
	x, _ := xml.Marshal(s) //生成xml

	//请求微信api
	gCurl.SetContentTypeXml()
	resStr := gCurl.RequestXml("POST", UNIFIEDORDER_URL, string(x))

	//解析返回xml
	err = xml.Unmarshal([]byte(resStr), &res)

	//记录日志
	gLog.Json("wechat.prepay.order.res", res)

	return
}
