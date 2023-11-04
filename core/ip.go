package core

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"regexp"
)

// AlterIP make some come true
// AlterIP 只选取部分实现
func AlterIP(ip string, formats []string) []string {
	var alteredIPs []string
	for _, format := range formats {
		standardIP := net.ParseIP(ip)
		switch format {
		case "1":
			alteredIPs = append(alteredIPs, standardIP.String())
		case "2":
			// 0-optimized dotted-decimal notation
			// the 0 value segments of an IP address can be ommitted (eg. 127.0.0.1 => 127.1)
			// regex for zeroes with dot 0000.
			var reZeroesWithDot = regexp.MustCompile(`(?m)[0]+\.`)
			// regex for .0000
			var reDotWithZeroes = regexp.MustCompile(`(?m)\.[0^]+$`)
			// suppress 0000.
			alteredIP := reZeroesWithDot.ReplaceAllString(standardIP.String(), "")
			// suppress .0000
			alteredIP = reDotWithZeroes.ReplaceAllString(alteredIP, "")
			alteredIPs = append(alteredIPs, alteredIP)
		case "3":
			// Octal notation (leading zeroes are required):
			// eg: 127.0.0.1 => 0177.0.0.01
			alteredIP := fmt.Sprintf("%#04o.%#o.%#o.%#o", standardIP[12], standardIP[13], standardIP[14], standardIP[15])
			alteredIPs = append(alteredIPs, alteredIP)
		case "4":
			alteredIPWithDots := fmt.Sprintf("%#x.%#x.%#x.%#x", standardIP[12], standardIP[13], standardIP[14], standardIP[15])
			alteredIPWithZeroX := fmt.Sprintf("0x%s", hex.EncodeToString(standardIP[12:]))
			alteredIPWithRandomPrefixHex, _ := RandomHex(5, standardIP[12:])
			alteredIPWithRandomPrefix := fmt.Sprintf("0x%s", alteredIPWithRandomPrefixHex)
			alteredIPs = append(alteredIPs, alteredIPWithDots, alteredIPWithZeroX, alteredIPWithRandomPrefix)
		case "5":
			// Decimal notation a.k.a dword notation
			// 127.0.0.1 => 2130706433
			bigIP, _, _ := IPToInteger(standardIP)
			alteredIPs = append(alteredIPs, bigIP.String())
		case "6":
			// Binary notation#
			// 127.0.0.1 => 01111111000000000000000000000001
			// converts to int
			bigIP, _, _ := IPToInteger(standardIP)
			// then to binary
			alteredIP := fmt.Sprintf("%b", bigIP)
			alteredIPs = append(alteredIPs, alteredIP)
		case "7":
			// Mixed notation
			// Ipv4 only
			alteredIP := fmt.Sprintf("%#x.%d.%#o.%#x", standardIP[12], standardIP[13], standardIP[14], standardIP[15])
			alteredIPs = append(alteredIPs, alteredIP)
		case "8":
			// URL-encoded IP address
			// 127.0.0.1 => %31%32%37%2E%30%2E%30%2E%31
			// ::1 => %3A%3A%31
			alteredIP := escape(ip)
			alteredIPs = append(alteredIPs, alteredIP)
		}
	}
	return alteredIPs
}

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

// IPToInteger converts an IP address to its integer representation.
// It supports both IPv4 as well as IPv6 addresses.
func IPToInteger(ip net.IP) (*big.Int, int, error) {
	// Binary operation, learn something
	// 二进制操作
	val := &big.Int{}
	val.SetBytes([]byte(ip))

	if len(ip) == net.IPv4len {
		return val, 32, nil //nolint
	} else if len(ip) == net.IPv6len {
		return val, 128, nil //nolint
	} else {
		return nil, 0, fmt.Errorf("unsupported address length %d", len(ip))
	}
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
