package lt

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/wshaman/ltntreader/tools/curl"
)

const (
	baseURL = "https://litnet.com"
	authURL = baseURL + "/auth/login?classic=1"
	libURL  = baseURL + "/account/library"
)

func NewLtWithIdentity(identity string) (l *Lt, err error) {
	l = &Lt{}
	l.client = curl.NewClient()
	if err = l.setIdentity(identity); err != nil {
		return nil, errors.Wrap(err, "failed to create Lt")
	}

	return l, nil
}

func NewLt(user, password string) (l *Lt, err error) {
	l = &Lt{}
	l.client = curl.NewClient()
	if err = l.login(user, password); err != nil {
		return nil, errors.Wrap(err, "failed to create Lt")
	}
	return l, nil
}

func (l *Lt) getLoginCSRF() (csrf string, err error) {
	var loginForm, csrfInput *goquery.Selection
	var ok bool
	var data *curl.Response

	if data, err = l.client.DoGet(authURL, nil, false); err != nil {
		return "", errors.Wrap(err, "failed to get login page")
	}
	if err = ioutil.WriteFile("/tmp/csrf.html", data.Body, 0700); err != nil {
		return "", errors.Wrap(err, "failed to write file")
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data.Body))
	if err != nil {
		return "", errors.Wrap(err, "failed to load parser login page")
	}
	if loginForm = doc.Find("form#w0"); loginForm == nil {
		return "", errors.New("no form#w0 found on given login page")

	}
	if csrfInput = loginForm.Find("input[name=\"_csrf\"]"); csrfInput == nil {
		return "", errors.New("no form csrf found on given login page")
	}
	if csrf, ok = csrfInput.Attr("value"); !ok {
		return "", errors.New("no csrf value is set on login form")
	}

	return csrf, nil
}

func (l *Lt) login(user, password string) (err error) {

	csrf, err := l.getLoginCSRF()
	if err != nil {
		return errors.Wrap(err, "login failed")
	}
	frm := url.Values{}
	frm.Add("LoginForm[login]", user)
	frm.Add("LoginForm[password]", password)
	frm.Add("_csrf", csrf)
	frm.Add("register-button", "")

	headers := map[string]string{
		":authority":                "litnet.com",
		":method":                   "POST",
		":path":                     "/auth/login?classic=1",
		":scheme":                   "https",
		"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",
		"accept-encoding":           "gzip, deflate, br",
		"accept-language":           "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
		"cache-control":             "max-age=0",
		"origin":                    "https://litnet.com",
		"referer":                   authURL,
		"upgrade-insecure-requests": "1",
	}
	resp, err := l.client.DoPost(authURL, strings.NewReader(frm.Encode()), headers)
	if err != nil {
		return errors.Wrap(err, "login failed")
	}
	if err = ioutil.WriteFile("/tmp/login.html", resp.Body, 0700); err != nil {
		return errors.Wrap(err, "failed to write file")
	}
	l.client.SaveCookies(resp.Resp.Cookies())
	return nil
}

func (l *Lt) setIdentity(identity string) (err error) {
	csrf, err := l.getLoginCSRF()
	if err != nil {
		return errors.Wrap(err, "failed to get csrf")
	}
	cookies := []*http.Cookie{{
		Name:     "_identity",
		Value:    identity,
		Domain:   "litnet.com",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: false,
		SameSite: 0,
		Unparsed: nil,
	},
		{
			Name:     "_csrf",
			Value:    csrf,
			Domain:   "litnet.com",
			MaxAge:   0,
			Secure:   false,
			HttpOnly: false,
			SameSite: 0,
			Unparsed: nil,
		},
	}
	l.client.SaveCookies(cookies)
	return nil
}
