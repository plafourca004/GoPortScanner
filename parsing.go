//Package main provides the parsng functions and the scanning functions
package main

import (
	"net"
	"regexp"
	"strconv"
	"strings"
)

//Get a list of ports from the parameters passed in the command line.
//You can either put your ports separated by a coma : 20,21,22,23,24,25.
//Or you can put your ports using a dash : 20-25
func getPortsList(ports string) []int {
	var portsList []int
	portsListString := strings.Split(ports, ",")

	for _, port := range portsListString {
		if strings.Contains(port, "-") { //number-number case
			var rangePorts [2]int
			rangePorts[0], _ = strconv.Atoi(strings.Split(port, "-")[0])
			rangePorts[1], _ = strconv.Atoi(strings.Split(port, "-")[1])
			for i := rangePorts[0]; i <= rangePorts[1]; i++ {
				portsList = append(portsList, i)
			}
		} else {
			portNb, _ := strconv.Atoi(port)
			portsList = append(portsList, portNb)
		}
	}

	return portsList
}

//Get a list of ip adresses from the parameters passed in the command line
//You can either put a single ip like 127.0.0.1 or a whole network like 192.168.1.0/30
func getIpList(ips []string) []string {
	var ipList []string

	regexCIDR, _ := regexp.Compile(`^([0-9]{1,3}\.){3}[0-9]{1,3}(\/([0-9]|[1-2][0-9]|3[0-2]))$`)

	for _, ip := range ips {
		if regexCIDR.MatchString(ip) { //Is a CIDR network adress
			ip4, ipNetwork, _ := net.ParseCIDR(ip)
			var ipListForThis []string
			for currentIp := ip4; ipNetwork.Contains(currentIp); nextIp(currentIp) {
				ipListForThis = append(ipListForThis, currentIp.String())
			}
			ipListForThis = ipListForThis[1:]                    //Remove network adress
			ipListForThis = ipListForThis[:len(ipListForThis)-1] //Remove broadcast adress

			ipList = append(ipList, ipListForThis...)
		} else {
			ipList = append(ipList, ip)
		}
	}

	return ipList
}

func nextIp(ip net.IP) net.IP {
	ip = ip.To4()
	if ip[3] == 255 {
		ip[2]++
	}
	ip[3]++

	return ip
}
