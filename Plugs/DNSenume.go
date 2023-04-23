package Plugs

import (
	"fmt"
	"github.com/wintrysec/FindSubs/Common"
	"net"
	"sync"
)

func DnsEnume(domain string, wg *sync.WaitGroup, subdomains chan<- string) {
	defer wg.Done()
	var wgdns sync.WaitGroup

	//判断是否存在泛解析
	if IsBad(domain) {
		//存在泛解析的处理流程
		DnsBadFlag, _ := net.LookupCNAME("ddddvvv000." + domain)
		fmt.Printf("泛解析对比标志为: %s\n", DnsBadFlag)
		for _, sub := range Common.Dictionary {
			wgdns.Add(1)
			subdomain := sub + "." + domain
			go func(subdomain string) {
				defer wgdns.Done()
				cnmae, err := net.LookupCNAME(subdomain)
				if err == nil {
					if cnmae == DnsBadFlag {
						//泛解析直接跳过
					} else {
						subdomains <- subdomain
						msg := fmt.Sprintf("[DNS_Enume Found] %s -> %v\n", subdomain, cnmae)
						fmt.Printf(msg)
					}
				}
			}(subdomain)
		}
	} else {
		//不存在泛解析的处理流程
		for _, sub := range Common.Dictionary {
			wgdns.Add(1)
			subdomain := sub + "." + domain
			go func(subdomain string) {
				defer wgdns.Done()
				ips, err := net.LookupIP(subdomain)
				if err == nil {
					subdomains <- subdomain
					msg := fmt.Sprintf("[DNS_Enume Found] %s -> %v\n", subdomain, ips)
					fmt.Printf(msg)
				} else {
					//fmt.Printf(fmt.Sprintf("不存在的域名:%s\n", subdomain))
				}
			}(subdomain)
		}
	}
	wgdns.Wait()

}

// 判断是否存在泛解析,存在则返回true
func IsBad(domain string) bool {
	subdomain := "ddddvvv000." + domain //肯定不存在的子域名
	ip, err := net.LookupIP(subdomain)
	if err == nil {
		fmt.Println("存在泛解析", ip)
		return true
	} else {
		fmt.Println("不存在泛解析")
		return false
	}
}
