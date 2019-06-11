package wechat

import (
	"fmt"
	"github.com/lucky-lee/gutil/gStr"
	"time"
)

type TwoSignJsapi struct {
	AppId     string `json:"appId" xml:"appId"`         //应用id
	NonceStr  string `json:"nonce_str" xml:"nonce_str"` //随机字符串 不长于32位
	Package   string `json:"package" xml:"package"`     //数据包
	SignType  string `json:"signType" xml:"signType"`   //签名方式
	TimeStamp string `json:"timeStamp" xml:"timeStamp"` //时间戳

	PaySign string `json:"paySign" xml:"-"` //
}

//生成 二次签名 配置
func newTwoSignJsapi(prepayId string) (res TwoSignJsapi) {
	res.AppId = AppId()
	res.TimeStamp = fmt.Sprintf("%d", time.Now().Unix())
	res.NonceStr = gStr.RandStr(32)
	res.Package = fmt.Sprintf("prepay_id=%s", prepayId)
	res.SignType = "MD5"
	res.PaySign = toSign(res)

	return
}
