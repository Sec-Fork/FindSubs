package Plugs

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Response struct {
	Results []struct {
		Data  []string `json:"parsed.names"`
		Data1 []string `json:"parsed.extensions.subject_alt_name.dns_names"`
	} `json:"results"`
}

// 去重函数
func unique(slice []string, domain string) []string {
	m := make(map[string]bool)
	for _, s := range slice {
		m[s] = true
	}
	var result []string
	for key, _ := range m {
		if strings.Contains(key, domain) {
			result = append(result, key)
		}

	}
	return result
}

func QueryCensys(apiID, apiSecret, domain string, wg *sync.WaitGroup, subdomains chan<- string) (err error) {
	defer wg.Done()

	//设置请求参数
	apiURL := "https://search.censys.io/api/v1/search/certificates"
	data := map[string]interface{}{
		"query":   domain, //替换为域名
		"fields":  []string{"parsed.names", "parsed.extensions.subject_alt_name.dns_names"},
		"flatten": true,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	req.SetBasicAuth(apiID, apiSecret)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	//设置http客户端参数
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //忽略https验证
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(12) * time.Second, //设置超时连接
	}

	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("[Censys-err] Failed to query %s", domain))
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("[Censys-err] Can't reading response body %s", domain))
		return err
	}
	//fmt.Println(string(body))

	// 定义一个正则表达式，用于匹配符合子域名规则的字符串
	pattern := `([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-zA-Z]{2,}`

	// 编译正则表达式
	reg := regexp.MustCompile(pattern)

	// 查找所有符合子域名规则的字符串
	subdomainss := reg.FindAllString(string(body), -1)

	// 去重
	subdomainss = unique(subdomainss, domain)

	// 输出结果
	for _, dom := range subdomainss {
		subdomains <- dom
		fmt.Printf("[Censys Found] %s\n", dom)
	}
	return nil
}
