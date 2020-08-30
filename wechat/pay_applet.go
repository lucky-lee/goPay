package wechat

import (
	"fmt"
	"github.com/lucky-lee/gutil/gStr"
	"time"
)

type TwoSignApplet struct {
	AppId     string `json:"app_id" xml:"appId"`        //应用id
	NonceStr  string `json:"nonce_str" xml:"nonceStr"`  //随机字符串 不长于32位
	Package   string `json:"package" xml:"package"`     //数据包
	SignType  string `json:"sign_type" xml:"signType"`  //签名方式
	Timestamp string `json:"timestamp" xml:"timeStamp"` //时间戳

	PaySign string `json:"pay_sign" xml:"-"` //
}

//生成 二次签名 配置
func newTwoSignApplet(orderId, prepayId string) (res TwoSignApplet) {
	res.AppId = AppId()
	res.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	res.NonceStr = gStr.RandStr(32)
	res.Package = fmt.Sprintf("prepay_id=%s", prepayId)
	res.SignType = "MD5"
	res.PaySign = toSign(res)

	return
}
