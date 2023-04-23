package Common

import (
	"flag"
)

func Flag(Info *HostInfo) {
	flag.StringVar(&Info.Domain, "domain", "", "The root domain of target")
	flag.StringVar(&Info.DomPath, "domfile", "./domains.txt", "The root domains file of target")
	flag.BoolVar(&Info.Savelog, "log", false, "Save log to ./domlog.txt")
	flag.BoolVar(&Info.SaveIP, "ips", false, "Save ips to ./ips.txt")
	flag.Parse()
}
