package wechat

import (
	"encoding/xml"
	"fmt"
	"git.dripcar.cn/golang/common/pubConstant/pubBean"
	"github.com/lucky-lee/gutil/gCurl"
	"github.com/lucky-lee/gutil/gLog"
	"github.com/lucky-lee/gutil/gStr"
	"reflect"
	"strings"
)

//支付配置
type ConfigPay struct {
	AppId      string `xml:"appid"`            //应用id
	Body       string `xml:"body"`             //商品描述
	MchId      string `xml:"mch_id"`           //商户号
	NonceStr   string `xml:"nonce_str"`        //随机字符串
	NotifyUrl  string `xml:"notify_url"`       //通知地址
	OutTradeNo string `xml:"out_trade_no"`     //商户订单号
	ClientIp   string `xml:"spbill_create_ip"` //终端IP
	TotalFee   string `xml:"total_fee"`        //总金额
	TradeType  string `xml:"trade_type"`       //交易类型
	Sign       string `xml:"sign"`             //签名
}

//设置-商品描述
func (b *ConfigPay) SetBody(s string) *ConfigPay {
	b.Body = s
	return b
}

//设置-商户订单号
func (b *ConfigPay) SetOutTradeNo(s string) *ConfigPay {
	b.OutTradeNo = s
	return b
}

//设置-总金额
func (b *ConfigPay) SetTotalFree(f float64) *ConfigPay {
	b.TotalFee = fmt.Sprintf("%.0f", f*100)
	return b
}

//设置-通知地址
func (b *ConfigPay) SetNotifyUrl(s string) *ConfigPay {
	b.NotifyUrl = s
	return b
}

//设置-终端ip
func (b *ConfigPay) SetClientIp(s string) *ConfigPay {
	b.ClientIp = s
	return b
}

//生成-预付订单
func (b ConfigPay) PrePayOrder() interface{} {
	b.Sign = toSign(b)     //生成签名字符串
	x, _ := xml.Marshal(b) //生成xml

	//请求微信api
	gCurl.SetContentTypeXml()
	resStr := gCurl.RequestXml("POST", UNIFIEDORDER_URL, string(x))

	//解析返回xml
	var beanRes ResPay
	err := xml.Unmarshal([]byte(resStr), &beanRes)

	if err != nil { //解析错误
		return pubBean.BeanEmpty{}
	}

	//记录日志
	gLog.Json("wechatPrePayOrderResponse", beanRes)

	//判断是否成功调用微信预支付接口
	if beanRes.ReturnCode == CODE_SUCCESS && beanRes.ResultCode == CODE_SUCCESS {
		return newTwoSignApp(beanRes.PrepayId)
	} else {
		return pubBean.BeanEmpty{}
	}
}

//微信支付 返回结果
type ResPay struct {
	ReturnCode string `xml:"return_code"` //
	ResultCode string `xml:"result_code"` //
	ReturnMsg  string `xml:"return_msg"`  //
	AppId      string `xml:"appid"`       //
	MchId      string `xml:"mchid"`       //
	NonceStr   string `xml:"nonce_str"`   //
	Sign       string `xml:"sign"`        //
	PrepayId   string `xml:"prepay_id"`   //
	TradeType  string `xml:"trade_type"`  //
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

		if key != "sign" && key != "packageval" { //过滤sign packageval
			s := fmt.Sprintf("%s=%s", key, val.String())
			arr = append(arr, s)
		}

	}

	return strings.Join(arr, "&")
}
