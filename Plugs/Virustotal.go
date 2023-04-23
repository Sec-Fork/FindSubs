package Plugs

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type response struct {
	Subdomains []string `json:"subdomains"`
}

// 查询virustotal API获取子域名
func QueryVirustotal(apiKey, domain string, wg *sync.WaitGroup, subdomains chan<- string) (err error) {
	defer wg.Done()

	//设置请求参数
	url := fmt.Sprintf("https://www.virustotal.com/vtapi/v2/domain/report?apikey=%s&domain=%s", apiKey, domain)
	req, _ := http.NewRequest("GET", url, nil)

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
		err = errors.New(fmt.Sprintf("[Virustotal-err] Failed to query %s", domain))
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("[Virustotal-err] Bad Status %s %s", domain, resp.Status))
		return err
	}

	var data response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		err = errors.New(fmt.Sprintf("[Virustotal-err] Json Decode is failed %s", domain))
		return err
	}

	for _, subdomain := range data.Subdomains {
		subdomains <- subdomain
		fmt.Printf("[Virustotal Found] %s\n", subdomain)
	}
	return nil
}
