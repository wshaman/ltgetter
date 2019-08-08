package curl

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/wshaman/ltntreader/tools"
	//"github.com/PuerkitoBio/goquery"
)

type Client struct {
	c       *http.Client
	cookies []*http.Cookie
}

type Response struct {
	Body []byte
	Resp *http.Response
}

var WaitFor = 8 * time.Second

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36"

func NewClient() *Client {
	return &Client{
		c: &http.Client{
			Timeout: 30 * time.Second,
		},
		cookies: nil,
	}
}

func newResp(httpResp *http.Response) (r *Response, err error) {
	r = &Response{}
	var b []byte
	if b, err = ioutil.ReadAll(httpResp.Body); err != nil {
		return nil, errors.Wrap(err, "failed to read body")
	}
	r.Body = b
	r.Resp = httpResp
	return r, nil
}

func (c *Client) DoGet(url string, headers map[string]string, saveCookie bool) (r *Response, err error) {
	tOut := tools.RndTime(WaitFor, 3*time.Second)
	fmt.Printf("waiting for %5d seconds ... ", tOut/time.Second)
	time.Sleep(tOut)
	fmt.Println("done")
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	response, err := c.do(request)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get url %s", url))
	}
	defer response.Body.Close()
	if saveCookie {
		c.SaveCookies(response.Cookies())
	}
	resp, err := newResp(response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wrap response")
	}
	savePage(request, resp.Body)
	return resp, err
}

func (c *Client) DoPost(url string, data io.Reader, headers map[string]string) (r *Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to post request")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := c.do(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to post url %s", url))
	}
	defer response.Body.Close()

	return newResp(response)
}

func (c *Client) SaveCookies(cookies []*http.Cookie) {
	c.cookies = cookies
}

func (c *Client) do(r *http.Request) (*http.Response, error) {
	for _, v := range c.cookies {
		r.AddCookie(v)
	}
	r.Header.Set("User-Agent", userAgent)
	return c.c.Do(r)
}

func savePage(rq *http.Request, body []byte) {
	_f := strings.Replace(rq.URL.String(), ":", "", -1)
	_f = strings.Replace(_f, "/", "_", -1)
	fn := path.Join(os.TempDir(), "ltnt")
	if err := os.MkdirAll(fn, 0722); err != nil {
		fmt.Println(err.Error())
	}
	fn = path.Join(fn, _f+".html")
	if err := ioutil.WriteFile(fn, body, 0644); err != nil {
		fmt.Println(err.Error())
	}
}
