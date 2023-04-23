package Plugs

import (
	"context"
	"fmt"
	"github.com/wintrysec/FindSubs/Common"
	"net"
	"sync"
	"time"
)

//type respofisp struct {
//	Country   string `json:"country"`
//	ShortName string `json:"short_name"`
//	City      string `json:"city"`
//	Isp       string `json:"isp"`
//	Net       string `json:"net"`
//}

// CollectSubdomains 收集子域名并去重
func CollectSubdomains(domains []string, apiKey, apiID, apiSecret string) []string {
	var subdomains []string
	subdomainMap := make(map[string]bool)
	subdomainsChan := make(chan string)

	var wg sync.WaitGroup
	for _, domain := range domains {
		wg.Add(3)
		go func(domain string) {
			QueryVirustotal(apiKey, domain, &wg, subdomainsChan)
		}(domain)

		go func(domain string) {
			QueryCensys(apiID, apiSecret, domain, &wg, subdomainsChan)
		}(domain)

		go func(domain string) {
			DnsEnume(domain, &wg, subdomainsChan)
		}(domain)
	}

	go func() {
		for subdomain := range subdomainsChan {
			subdomainMap[subdomain] = true
		}
	}()

	wg.Wait()
	close(subdomainsChan)

	for subdomain := range subdomainMap {
		subdomains = append(subdomains, subdomain)
	}

	return subdomains
}

// DNS解析和CDN识别
func CDNLook(subdomain string, SaveIP bool) {
	//自定义nameserver
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 10 * time.Second,
			}
			return d.DialContext(ctx, "udp", "8.8.8.8:53")
		},
	}
	//解析DNS为HOST
	ips, err := r.LookupHost(context.Background(), subdomain)
	ipaddr := "CDN"
	if err == nil {
		//isp := "-"

		if len(ips) == 1 {
			//判断CDN
			ipaddr = ips[0]
			if SaveIP {
				Common.SaveIPS(ipaddr)
			}
			//isp, err = QueryISP(ipaddr)
			//if err != nil {
			//	fmt.Println(err)
			//}
		}
		msg := fmt.Sprintf("%s %v", subdomain, ipaddr)
		fmt.Println(msg)
	} else {
		//无法解析出IP
		msg := fmt.Sprintf("%s", subdomain)
		fmt.Println(msg)
	}
}

//// 查询IP归属
//func QueryISP(ip string) (isp string, err error) {
//	//请求参数设置
//	url := fmt.Sprintf("https://ip.useragentinfo.com/json?ip=%s", ip)
//	req, _ := http.NewRequest("GET", url, nil)
//
//	//设置http客户端参数
//	tr := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //忽略https验证
//	}
//	client := &http.Client{
//		Transport: tr,
//		Timeout:   time.Duration(5) * time.Second, //设置超时连接
//	}
//
//	//发送HTTP请求
//	resp, err := client.Do(req)
//	if err != nil {
//		err = errors.New(fmt.Sprintf("[QueryISP-err] HTTP can‘t Connect %s", ip))
//		return "", err
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		err = errors.New(fmt.Sprintf("[QueryISP-err] HTTP Status is not ok %s", ip))
//		return "", err
//	}
//	defer resp.Body.Close()
//
//	//响应解析
//	var data respofisp
//	err = json.NewDecoder(resp.Body).Decode(&data)
//
//	//格式处理
//	var Country = ""
//	if data.City != "" {
//		Country = data.City
//	} else if data.ShortName != "" {
//		Country = data.ShortName
//	} else {
//		Country = data.Country
//	}
//	base_isp := data.Isp + data.Net
//	isp = fmt.Sprintf("%s(%s)", base_isp, Country)
//	return isp, nil
//}
