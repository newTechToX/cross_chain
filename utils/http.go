package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/net/context/ctxhttp"
)

type HttpOption struct {
	Method      string
	Host        string
	Url         *url.URL
	Header      map[string]string
	RequestBody interface{}
	Response    interface{}
	Proxy       string
}

func (ho *HttpOption) Send(ctx context.Context) error {
	log.Debug("http option send", "method", ho.Method, "url", ho.Url, "proxy", ho.Proxy)
	if ho.Url == nil {
		return fmt.Errorf("no url specificed")
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(ho.RequestBody); err != nil {
		return err
	}
	req, err := http.NewRequest(ho.Method, ho.Url.String(), &buf)
	if err != nil {
		return err
	}

	if ho.Host != "" {
		req.Host = ho.Host
	}
	if ho.Header != nil {
		for k, v := range ho.Header {
			req.Header.Set(k, v)
		}
	}
	client := &http.Client{}
	if ho.Proxy != "" {
		proxy, err := url.Parse(ho.Proxy)
		if err != nil {
			return err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}
	resp, err := ctxhttp.Do(ctx, client, req)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, ho.Response); err != nil {
		return err
	}
	_ = resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("StatusCode: %d", resp.StatusCode)
	}
	return nil
}

func IsHttps(uri string) bool {
	return strings.Index(strings.ToLower(uri), "https://") == 0
}

func IsHttp(uri string) bool {
	return strings.Index(strings.ToLower(uri), "http://") == 0
}

func HttpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func HttpGetWithProxy(uri, proxyUrl string) ([]byte, error) {
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Transport: &http.Transport{
			// 设置代理
			Proxy: http.ProxyURL(proxy),
		},
	}
	resp, err := client.Get(uri)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func HttpGetObjectWithProxy(url, proxy string, dest any) error {
	if proxy == "" {
		return HttpGetObject(url, dest)
	}
	data, err := HttpGetWithProxy(url, proxy)
	// fmt.Println(string(data))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("%v, raw data: %v", err, string(data))
	}

	return nil
}

func HttpGetObject(url string, dest any) error {
	data, err := HttpGet(url)
	if err != nil {
		return err
	}
	// fmt.Println("======")
	// fmt.Println(string(data))
	// fmt.Println("======")

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("%v, raw data: %v", err, string(data))
	}

	return nil
}

// dns resolve, tcp connect refuse, timeout
func IsNetError(err error) bool {
	if _, ok := err.(net.Error); ok {
		return true
	}
	return false
}
