package core

import (
	"math/big"
	"net"
)

func IPRange(startIPStr, endIPStr string) []string {
	var ipList []string
	startIPInt := ipToUInt32(startIPStr)
	endIPInt := ipToUInt32(endIPStr)
	for currIPInt := new(big.Int).Set(startIPInt); currIPInt.Cmp(endIPInt) <= 0; incIP(currIPInt) {
		ip := intToIP(currIPInt)
		ipList = append(ipList, ip)
	}

	return ipList
}

func intToIP(ipInt *big.Int) string {
	ipBytes := ipInt.Bytes()
	if len(ipBytes) < 4 {
		// 补齐 4 个字节
		padBytes := make([]byte, 4-len(ipBytes))
		ipBytes = append(padBytes, ipBytes...)
	}
	return net.IP(ipBytes).String()
}

func ipToUInt32(ipStr string) *big.Int {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}

	// 将 net.IP 转换为 4 字节的大整数
	ipInt := new(big.Int)
	ipInt.SetBytes(ip.To4())
	return ipInt
}

func incIP(ipInt *big.Int) {
	ipInt.Add(ipInt, big.NewInt(1))
}

//func check
//func main() {
//startIPStr := "192.168.1.10"
//endIPStr := "192.168.10.15"
//
//startIPInt := ipToUInt32(startIPStr)
//endIPInt := ipToUInt32(endIPStr)
//
//if startIPInt != nil && endIPInt != nil {
//	ipList := ipRange(startIPInt, endIPInt)
//	for _, ip := range ipList {
//		fmt.Println(ip)
//	}
//} else {
//	fmt.Println("无效的IP地址")
//}
//}
