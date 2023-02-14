package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type client struct {
	Config *Config
	Cookie *http.Cookie
}

type Config struct {
	Headers    map[string]string //http header
	Timeout    int64             //请求超时时间
	CookieName string            //Cookie名称
	Username   string            //basicAuth用户名
	Password   string            //basicAuth密码
}

func NewClient(in Config) *client {
	return &client{
		Config: &Config{
			Headers:    in.Headers,
			Timeout:    in.Timeout,
			CookieName: in.CookieName,
			Username:   in.Username,
			Password:   in.Password,
		},
	}
}

type VarServer struct {
	Url       string            //请求地址
	BasicAuth bool              //是否需要Basic认真
	Data      interface{}       //请求体
	IsPrintln bool              //是否打印结果
	Fields    map[string]string //上传文件时的参数、值
	Body      io.Reader         //文件内容
}

func (e *client) Post(in VarServer) ([]byte, error) {
	var (
		inData, result []byte
		err            error
	)
	if inData, err = json.Marshal(in.Data); err != nil {
		return result, err
	}
	if result, err = e._request(http.MethodPost, in.Url, in.BasicAuth, bytes.NewReader(inData)); err != nil {
		return result, err
	}
	if !json.Valid(result) {
		return result, errors.New("result is not standard JSON")
	}
	if in.IsPrintln {
		log.Println(fmt.Sprintf("本次[%s]请求[%s]返回信息是: %s", http.MethodPost, in.Url, string(result)))
	}
	return result, err
}

func (e *client) Get(in VarServer) ([]byte, error) {
	var (
		result []byte
		err    error
	)
	if result, err = e._request(http.MethodGet, in.Url, in.BasicAuth, nil); err != nil {
		return result, err
	}
	if !json.Valid(result) {
		return result, errors.New("result is not standard JSON")
	}
	if in.IsPrintln {
		log.Println(fmt.Sprintf("本次[%s]请求[%s]返回信息是: %s", http.MethodGet, in.Url, string(result)))
	}
	return result, err
}

func (e *client) Upload(in VarServer) ([]byte, error) {
	var (
		result []byte
		err    error
		body   = &bytes.Buffer{}
		writer = multipart.NewWriter(body)
		fw     io.Writer
	)
	if fw, err = writer.CreateFormFile("file", in.Fields["filename"]); err != nil {
		return result, err
	}
	if _, err = io.Copy(fw, in.Body); err != nil {
		return result, err
	}
	for key, value := range in.Fields {
		if err = writer.WriteField(key, value); err != nil {
			return result, err
		}
	}
	if err = writer.Close(); err != nil {
		return result, err
	}
	e.Config.Headers["Content-Type"] = writer.FormDataContentType()
	if result, err = e._request(http.MethodPost, in.Url, in.BasicAuth, body); err != nil {
		return result, err
	}
	if in.IsPrintln {
		log.Println(fmt.Sprintf("本次[%s]请求[%s]返回信息是: %s", http.MethodPost, in.Url, string(result)))
	}
	if !json.Valid(result) {
		return result, errors.New("result is not standard JSON")
	}
	return result, err
}

func (e *client) _request(method, url string, basicAuth bool, body io.Reader) (result []byte, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)

	httpClient := http.Client{
		Timeout: time.Duration(e.Config.Timeout) * time.Second,
	}
	if req, err = http.NewRequest(method, url, body); err != nil {
		return
	}
	if len(e.Config.Headers) > 0 {
		for key, value := range e.Config.Headers {
			req.Header.Set(key, value)
		}
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	}
	if basicAuth {
		req.SetBasicAuth(e.Config.Username, e.Config.Password)
	}
	if resp, err = httpClient.Do(req); err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if e.Config.CookieName != "" && resp.Cookies() != nil {
		for _, v := range resp.Cookies() {
			if v.Name == e.Config.CookieName {
				e.Cookie = v
			}
		}
	}
	return io.ReadAll(resp.Body)
}
