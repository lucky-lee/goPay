package wechat

//异步通知
//微信支付-异步通知实体
type Notify struct {
	AppId         string `json:"appid" xml:"appid"`                   //应用id
	MchId         string `json:"mch_id" xml:"mch_id"`                 //商户号
	BankType      string `json:"bank_type" xml:"bank_type"`           //付款银行
	CashFee       int    `json:"cash_fee" xml:"cash_fee"`             //现金支付金额
	FeeType       string `json:"fee_type" xml:"fee_type"`             //货币种类
	IsSubscribe   string `json:"is_subscribe" xml:"is_subscribe"`     //是否关注公众号 (Y是,N否)
	NonceStr      string `json:"nonce_str" xml:"nonce_str"`           //随机字符串
	OpenId        string `json:"openid" xml:"openid"`                 //用户标识
	OutTradeNo    string `json:"out_trade_no" xml:"out_trade_no"`     //商户订单号
	ResultCode    string `json:"result_code" xml:"result_code"`       //业务结果
	ReturnCode    string `json:"return_code" xml:"return_code"`       //返回状态码
	Sign          string `json:"sign" xml:"sign"`                     //签名
	TimeEnd       string `json:"time_end" xml:"time_end"`             //支付完成时间
	TotalFee      int    `json:"total_fee" xml:"total_fee"`           //总金额
	TradeType     string `json:"trade_type" xml:"trade_type"`         //交易类型
	TransactionId string `json:"transaction_id" xml:"transaction_id"` //微信支付订单号
}

//是否成功
func (s *Notify) IsSuccess() (status bool, msg string) {
	sign := toSign(s)

	if sign != s.Sign {
		msg = "签名验证失败"
		return
	}

	if s.ResultCode == "SUCCESS" && s.ReturnCode == "SUCCESS" {
		status = true
		return
	}

	msg = "验证失败"
	return
}
