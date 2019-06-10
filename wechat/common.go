package wechat

import (
	"encoding/xml"
	"fmt"
	"github.com/lucky-lee/gutil/gCurl"
	"github.com/lucky-lee/gutil/gLog"
	"github.com/lucky-lee/gutil/gStr"
	"reflect"
	"strings"
)

//支付配置 -按照ascii 排序
type ConfigPay struct {
	AppId      string `xml:"appid"`        //必传-应用id
	Attach     string `xml:"attach"`       //非必传-附加数据
	Body       string `xml:"body"`         //必传-商品描述
	Detail     string `xml:"detail"`       //非必传-商品详情
	DeviceInfo string `xml:"device_info"`  //非必传-设备号
	FeeType    string `xml:"fee_type"`     //非必传-标价币种 符合ISO 4217标准的三位字母代码，默认人民币：CNY
	GoodsTag   string `xml:"goods_tag"`    //非必传-订单优惠标记，使用代金券或立减优惠功能时需要的参数
	LimitPay   string `xml:"limit_pay"`    //非必传-上传此参数no_credit--可限制用户不能使用信用卡支付
	MchId      string `xml:"mch_id"`       //必传-商户号
	NonceStr   string `xml:"nonce_str"`    //必传-随机字符串
	NotifyUrl  string `xml:"notify_url"`   //必传-通知地址
	Openid     string `json:"openid"`      //非必传-trade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识
	OutTradeNo string `xml:"out_trade_no"` //必传-商户订单号

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

//生成-预付订单

func (self ConfigPay) PrePayOrder() (res interface{}) {
	res = ResEmpty{}
	self.Sign = toSign(self)  //生成签名字符串
	x, _ := xml.Marshal(self) //生成xml

	//请求微信api
	gCurl.SetContentTypeXml()
	resStr := gCurl.RequestXml("POST", UNIFIEDORDER_URL, string(x))

	//解析返回xml
	var resPay ResPay
	err := xml.Unmarshal([]byte(resStr), &resPay)

	if err != nil { //解析错误
		return ResEmpty{}
	}

	//记录日志
	gLog.Json("wechat.prepay.order.res", resPay)

	//判断是否成功调用微信预支付接口
	if resPay.ReturnCode == CODE_SUCCESS && resPay.ResultCode == CODE_SUCCESS {
		switch self.TradeType {
		case TYPE_APP:
			res = newTwoSignApp(resPay.PrepayId)
		case TYPE_JSAPI:
			res = newTwoSignJsapi(resPay.PrepayId)
		}

	}
	return
}

//空结构体
type ResEmpty struct {
}

//微信支付 返回结果
type ResPay struct {
	ReturnCode string `xml:"return_code"` //返回状态码
	ReturnMsg  string `xml:"return_msg"`  //返回信息

	//以下字段在return_code为SUCCESS的时候有返回
	AppId      string `xml:"appid"`        //应用id
	MchId      string `xml:"mch_id"`       //商户号
	DeviceInfo string `xml:"device_info"`  //设备号
	NonceStr   string `xml:"nonce_str"`    //随机字符串
	Sign       string `xml:"sign"`         //签名
	ResultCode string `xml:"result_code"`  //业务结果
	ErrCode    string `xml:"err_code"`     //错误代码
	ErrCodeDes string `xml:"err_code_des"` //错误代码描述

	//以下字段在return_code 和result_code都为SUCCESS的时候有返回
	TradeType string `xml:"trade_type"` //交易类型
	PrepayId  string `xml:"prepay_id"`  //预支付交易会话标识
	//trade_type=NATIVE时有返回，此url用于生成支付二维码，然后提供给用户进行扫码支付
	CodeUrl string `xml:"code_url"`
}

//生成签名字符串
//获取签名
// 签名生成的通用步骤如下：
// 第一步，设所有发送或者接收到的数据为集合M，将集合M内非空参数值的参数按照参数名ASCII码从小到大排序（字典序），使用URL键值对的格式（即key1=value1&key2=value2…）拼接成字符串stringA。
// 特别注意以下重要规则：
// ◆ 参数名ASCII码从小到大排序（字典序）；
// ◆ 如果参数的值为空不参与签名；
// ◆ 参数名区分大小写；
// ◆ 验证调用返回或微信主动通知签名时，传送的sign参数不参与签名，将生成的签名与该sign值作校验。
// ◆ 微信接口可能增加字段，验证签名时必须支持增加的扩展字段
// 第二步，在stringA最后拼接上key得到stringSignTemp字符串，并对stringSignTemp进行MD5运算，再将得到的字符串所有字符转换为大写，得到sign值signValue。
// key设置路径：微信商户平台(pay.weixin.qq.com)-->账户设置-->API安全-->密钥设置
func toSign(b interface{}) string {
	s := toSignStr(b)
	s += "&key=" + ApiKey()
	return strings.ToUpper(gStr.Md5(s))
}

//生成-签名字符串
func toSignStr(p interface{}) (signStr string) {
	var arr []string

	v := reflect.ValueOf(p)
	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Tag.Get("xml")
		val := v.Field(i)

		if key != "-" && val.String() != "" { //过滤sign packageval
			s := fmt.Sprintf("%s=%s", key, val.String())
			arr = append(arr, s)
		}

	}

	return strings.Join(arr, "&")
}
