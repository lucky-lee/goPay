package wechat

import (
	"fmt"
	"github.com/lucky-lee/gutil/gStr"
	"reflect"
	"sort"
	"strings"
)

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
	keys := make([]string, 0)
	val := make(map[string]string)
	signs := make([]string, 0)
	rVal := reflect.ValueOf(p)

	for i := 0; i < rVal.NumField(); i++ {
		field := rVal.Type().Field(i)
		k := field.Tag.Get("xml")
		v := rVal.Field(i)

		if k != "-" && v.String() != "" { //过滤sign packageval
			val[k] = v.String()
			keys = append(keys, k)
		}
	}

	//按照ASCII从小到大排序
	sort.Strings(keys)

	for _, v := range keys {
		sign := fmt.Sprintf("%s=%s", v, val[v])
		signs = append(signs, sign)
	}

	return strings.Join(signs, "&")
}
