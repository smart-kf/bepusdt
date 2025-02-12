package mail

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"regexp"
)

func validateEmail(email string) bool {
	// 正则表达式模式
	// 这个示例中的正则表达式并不是非常完整，仅供演示目的
	// 实际上，验证邮箱地址的正则表达式可能更复杂
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// 编译正则表达式
	regexpEmail := regexp.MustCompile(pattern)

	// 使用正则表达式验证邮箱
	return regexpEmail.MatchString(email)
}

type SendMailObject struct {
	To         string `json:"to"`
	From       string `json:"from"`
	HtmlString string `json:"html_string"`
}

func SendMail(ctx context.Context, object SendMailObject) error {
	if object.To == "" || !validateEmail(object.To) {
		return errors.New("参数错误")
	}

	var u = `https://api.elasticemail.com/v2/email/send?`
	var q = make(url.Values)
	q.Add("apikey", `6F894BB15DC6800B670AA84EB191AFD58E03758F5BFB3DC27CB837657926C8F55E29E25DC1DFC570674F465E2307E551`)
	q.Add("subject", "订单信息")
	q.Add("from", object.From)
	q.Add("to", object.To)
	q.Add("bodyHtml", object.HtmlString)
	var req, _ = http.NewRequest(http.MethodPost, u+q.Encode(), nil)
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp != nil {
		resp.Body.Close()
	}
	return nil
}
