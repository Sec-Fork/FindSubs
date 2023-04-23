package main

import (
	"fmt"
	"github.com/wintrysec/FindSubs/Common"
	"github.com/wintrysec/FindSubs/Plugs"
)

func main() {
	Common.Banner()
	//根域名列表
	var domains []string

	//解析命令选项
	var Info Common.HostInfo
	Common.Flag(&Info)
	if Info.Domain != "" {
		domains = append(domains, Info.Domain)
	} else {
		if Info.DomPath != "" {
			// 否则从文本文件中读取域名列表
			filepath := Info.DomPath
			domains, _ = Common.ReadDomainFromFile(filepath)
		} else {
			fmt.Println("Please input domain target")
		}
	}

	//分配收集任务
	Apikey, ApiID, ApiSecret := Common.GetApiKey()
	subdomains := Plugs.CollectSubdomains(domains, Apikey, ApiID, ApiSecret)

	//进入DNS解析和CDN识别流程 输出结果
	fmt.Printf("\n===================Total Found %d Subdomains===================\n", len(subdomains))
	for _, subdomain := range subdomains {
		if Info.Savelog {
			Common.Savelog(subdomain)
		}
		Plugs.CDNLook(subdomain, Info.SaveIP)
	}

}
