package wechat

import "os"

const (
	UNIFIEDORDER_URL = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	CODE_SUCCESS     = "SUCCESS"
	ENV_PREFIX       = "gopay.wechat." //env前缀
	ENV_APPID        = ENV_PREFIX + "appid"
	ENV_APIKEY       = ENV_PREFIX + "apikey"
	ENV_MCHID        = ENV_PREFIX + "mchid"

	//支付类型
	TYPE_APP    = "APP"
	TYPE_JSAPI  = "JSAPI"
	TYPE_NATIVE = "NATIVE"
	TYPE_H5     = "H5"

	//来源
	SOURCE_APP uint8 = iota + 1
	SOURCE_APPLET
)

//设置 appid
func SetAppId(val string) {
	os.Setenv(ENV_APPID, val)
}

//获取 appid
func AppId() string {
	return os.Getenv(ENV_APPID)
}

//32位api秘钥
func SetApiKey(val string) {
	os.Setenv(ENV_APIKEY, val)
}

//获取 apikey
func ApiKey() string {
	return os.Getenv(ENV_APIKEY)
}

//设置 商户秘钥
func SetMchId(val string) {
	os.Setenv(ENV_MCHID, val)
}

//获取 商户秘钥
func MchId() string {
	return os.Getenv(ENV_MCHID)
}
