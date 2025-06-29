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
)

// HttpClient HTTP客户端结构体
type HttpClient struct {
	ProxyOpen bool
	ProxyURL  string
}

// NewHttpClient 创建新的HTTP客户端实例
func NewHttpClient() *HttpClient {
	return &HttpClient{
		ProxyOpen: false,
		ProxyURL:  "",
	}
}

// NewHttpClientWithProxy 创建带代理配置的HTTP客户端实例
func NewHttpClientWithProxy(proxyOpen bool, proxyURL string) *HttpClient {
	return &HttpClient{
		ProxyOpen: proxyOpen,
		ProxyURL:  proxyURL,
	}
}

// SetProxy 设置代理配置
func (h *HttpClient) SetProxy(proxyOpen bool, proxyURL string) *HttpClient {
	h.ProxyOpen = proxyOpen
	h.ProxyURL = proxyURL
	return h
}

// Get 执行GET请求
func (h *HttpClient) Get(ctx context.Context, url string, header map[string]string, data interface{}, result interface{}, cookies ...map[string]string) error {
	client := g.Client().Timeout(60 * time.Second)
	if header != nil {
		client.SetHeaderMap(header)
	}

	if len(cookies) > 0 {
		client.Cookie(cookies[0])
	}

	// 设置代理
	if h.ProxyOpen && len(h.ProxyURL) > 0 {
		client.SetProxy(h.ProxyURL)
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

	if len(bytes) > 0 {
		err = gjson.Unmarshal(bytes, result)
		if err != nil {
			g.Log().Error(ctx, err)
			return err
		}
	}

	return nil
}

// Post 执行POST请求
func (h *HttpClient) Post(ctx context.Context, url string, header map[string]string, data, result interface{}) error {
	client := g.Client().Timeout(60 * time.Second)
	if header != nil {
		client.SetHeaderMap(header)
	}

	// 设置代理
	if h.ProxyOpen && len(h.ProxyURL) > 0 {
		client.SetProxy(h.ProxyURL)
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

	if len(bytes) > 0 {
		err = gjson.Unmarshal(bytes, result)
		if err != nil {
			g.Log().Error(ctx, err)
			return err
		}
	}

	return nil
}

// PostResult 执行POST请求并返回结果
func (h *HttpClient) PostResult(ctx context.Context, url string, header map[string]string, data, result interface{}) (res interface{}, err error) {
	client := g.Client().Timeout(60 * time.Second)
	if header != nil {
		client.SetHeaderMap(header)
	}

	if h.ProxyOpen && len(h.ProxyURL) > 0 {
		client.SetProxy(h.ProxyURL)
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

	if len(bytes) > 0 {
		err = gjson.Unmarshal(bytes, result)
		if err != nil {
			g.Log().Error(ctx, err)
			return
		}
	}
	res = result
	return
}

// DownloadFile 下载文件
func (h *HttpClient) DownloadFile(ctx context.Context, fileURL string, useProxy ...bool) []byte {
	g.Log().Debugf(ctx, "HttpDownloadFile fileURL: %s", fileURL)

	client := g.Client().Timeout(600 * time.Second)

	transport := &http.Transport{}

	if h.ProxyOpen && len(h.ProxyURL) > 0 && (len(useProxy) == 0 || useProxy[0]) {
		g.Log().Debugf(ctx, "HttpDownloadFile ProxyURL: %s", h.ProxyURL)

		proxyUrl, err := url.Parse(h.ProxyURL)
		if err != nil {
			g.Log().Error(ctx, err)
		}

		transport.Proxy = http.ProxyURL(proxyUrl)
		client.Transport = transport
	}

	return client.GetBytes(ctx, fileURL)
}

// GetProxy 获取代理函数
func (h *HttpClient) GetProxy(ctx context.Context) func(*http.Request) (*url.URL, error) {
	if h.ProxyOpen && len(h.ProxyURL) > 0 {
		g.Log().Debugf(ctx, "ProxyURL: %s", h.ProxyURL)

		proxyURL, err := url.Parse(h.ProxyURL)
		if err != nil {
			g.Log().Error(ctx, err)
			return nil
		}

		return http.ProxyURL(proxyURL)
	}

	return nil
}

// GetProxyTransport 获取代理传输器
func (h *HttpClient) GetProxyTransport(ctx context.Context) *http.Transport {
	transport := &http.Transport{}
	transport.Proxy = h.GetProxy(ctx)
	return transport
}
