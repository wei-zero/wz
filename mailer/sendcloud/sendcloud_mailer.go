package sendcloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/wei-zero/wz/mailer"
)

type SendcloudMailer struct {
	apiUser string
	apiKey  string
}

func NewMailer(apiUser, apiKey string) mailer.Mailer {
	return &SendcloudMailer{
		apiKey:  apiKey,
		apiUser: apiUser,
	}
}

func (s *SendcloudMailer) SendByTemplate(from *mailer.Email, to *mailer.Email, subject string, templateID string, data map[string]string) error {
	subs := make(map[string]any)
	for k, v := range data {
		subs["%"+k+"%"] = []string{v}
	}
	vars, _ := json.Marshal(map[string]any{
		"to":  []string{to.Address},
		"sub": subs,
	})

	params := url.Values{}
	params.Add("apiUser", s.apiUser)
	params.Add("apiKey", s.apiKey)
	params.Add("from", from.Address)
	params.Add("fromName", from.Name)
	params.Add("templateInvokeName", templateID)
	params.Add("subject", subject)
	params.Add("xsmtpapi", string(vars))
	params.Add("respEmailId", "true")

	//log.Println(params.Encode())
	res, err := http.PostForm("http://api.sendcloud.net/apiv2/mail/sendtemplate", params)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("send mail failed, code=%d, body: %s", res.StatusCode, string(b))
	}

	return nil
}
