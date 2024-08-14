package util

import (
	"context"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/util/gconv"

	"net/http"
	"net/url"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var ProxyOpen bool
var ProxyURL string

func init() {

	ctx := gctx.New()

	proxy_open, err := g.Cfg().Get(ctx, "http.proxy_open")
	if err != nil {
		g.Log().Error(ctx, err)
	}
	ProxyOpen = proxy_open.Bool()

	proxyUrl, err := g.Cfg().Get(ctx, "http.proxy_url")
	if err != nil {
		g.Log().Error(ctx, err)
	}

	ProxyURL = proxyUrl.String()
}

func HttpGet(ctx context.Context, url string, header map[string]string, data interface{}, result interface{}, cookies ...map[string]string) error {

	client := g.Client().Timeout(60 * time.Second)
	if header != nil {
		client.SetHeaderMap(header)
	}

	if len(cookies) > 0 {
		client.Cookie(cookies[0])
	}

	response, err := client.Get(ctx, url, data)
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}

	defer func(response *gclient.Response) {
		err = response.Close()
		if err != nil {
			g.Log().Error(ctx, err)
		}
	}(response)

	bytes := response.ReadAll()
	g.Log().Debugf(ctx, "HttpGet url: %s, header: %+v, data: %+v, response: %s", gconv.String(url), gconv.String(header), gconv.String(data), string(bytes))

	if bytes != nil && len(bytes) > 0 {
		err = gjson.Unmarshal(bytes, result)
		if err != nil {
			g.Log().Error(ctx, err)
			return err
		}
	}

	return nil
}

// HttpPost
//
//	@Description: post请求
//	@param ctx
//	@param url
//	@param header
//	@param data 发送的参数
//	@param result 返回的数据
//	@return error
func HttpPost(ctx context.Context, url string, header map[string]string, data, result interface{}) error {

	// g.Log().Debugf(ctx, "HttpPost url: %s, header: %+v, data: %+v", url, header, data)

	client := g.Client().Timeout(60 * time.Second)
	if header != nil {
		client.SetHeaderMap(header)
	}

	//设置代理
	if ProxyOpen && len(ProxyURL) > 0 {
		client.SetProxy(ProxyURL)
	}

	response, err := client.ContentJson().Post(ctx, url, data)
	if err != nil {
		g.Log().Error(ctx, err)
		return err
	}

	defer func(response *gclient.Response) {
		err = response.Close()
		if err != nil {
			g.Log().Error(ctx, err)
		}
	}(response)

	bytes := response.ReadAll()
	g.Log().Debugf(ctx, "HttpPost url: %s, header: %+v, data: %+v, response: %s", gconv.String(url), gconv.String(header), gconv.String(data), string(bytes))

	if bytes != nil && len(bytes) > 0 {
		err = gjson.Unmarshal(bytes, result)
		if err != nil {
			g.Log().Error(ctx, err)
			return err
		}
	}

	return nil
}

// 返回结果
func HttpPostResult(ctx context.Context, url string, header map[string]string, data, result interface{}) (res interface{}, err error) {

	client := g.Client().Timeout(60 * time.Second)
	if header != nil {
		client.SetHeaderMap(header)
	}

	if ProxyOpen && len(ProxyURL) > 0 {
		client.SetProxy(ProxyURL)
	}

	response, err := client.ContentJson().Post(ctx, url, data)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}

	defer func() {
		err = response.Close()
		if err != nil {
			g.Log().Error(ctx, err)
		}
	}()

	bytes := response.ReadAll()
	g.Log().Debugf(ctx, "HttpPost url: %s, header: %+v, data: %+v, response: %s", url, header, data, string(bytes))

	if bytes != nil && len(bytes) > 0 {
		err = gjson.Unmarshal(bytes, result)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	res = result
	return
}

func HttpDownloadFile(ctx context.Context, fileURL string, useProxy ...bool) []byte {

	g.Log().Debugf(ctx, "HttpDownloadFile fileURL: %s", fileURL)

	client := g.Client().Timeout(600 * time.Second)

	transport := &http.Transport{}

	if ProxyOpen && len(ProxyURL) > 0 && (len(useProxy) == 0 || useProxy[0]) {

		g.Log().Debugf(ctx, "HttpDownloadFile ProxyURL: %s", ProxyURL)

		proxyUrl, err := url.Parse(ProxyURL)
		if err != nil {
			g.Log().Error(ctx, err)
		}

		transport.Proxy = http.ProxyURL(proxyUrl)
		client.Transport = transport
	}

	return client.GetBytes(ctx, fileURL)
}

func GetProxy(ctx context.Context) func(*http.Request) (*url.URL, error) {

	var proxy func(*http.Request) (*url.URL, error)

	if ProxyOpen && len(ProxyURL) > 0 {

		g.Log().Debugf(ctx, "ProxyURL: %s", ProxyURL)

		proxyURL, err := url.Parse(ProxyURL)
		if err != nil {
			g.Log().Error(ctx, err)
			return nil
		}

		return http.ProxyURL(proxyURL)
	}

	return proxy
}

func GetProxyTransport(ctx context.Context) *http.Transport {

	transport := &http.Transport{}

	transport.Proxy = GetProxy(ctx)

	return transport
}
