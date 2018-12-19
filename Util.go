package gothic

/**
 * @desc Util
 * @author zhaojiangwei
 * @date 2018-12-18 10:35
 */

import (
	"net"
	"fmt"
	"time"
	"reflect"
	"os"
	"crypto/md5"
	"io"
	"encoding/binary"
	"net/http"
	"strings"
	"errors"
	"math/rand"
	"strconv"
)

const (
	XForwardedFor = "X-Forwarded-For"
	XRealIP       = "X-Real-IP"
)

//获取本机ip
func GetLocalIp() string {
	addrs, _ := net.InterfaceAddrs()
	var ip string = ""
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			ip = ipnet.IP.String()
			if ip != "127.0.0.1" {
			}
		}
	}
	return ip
}

//获取客户端ip
func GetClientIp(r *http.Request) string {
	//todo "X-Forwarded-For"
	socket := r.RemoteAddr
	if socket != "" {
		return strings.Split(r.RemoteAddr, ":")[0]
	} else {
		return ""
	}
}

//简单序列化成php的格式, 和已有的php系统交互的时候可能会用到
func SerializePhp(data map[string]interface{}) string {
	ret := fmt.Sprintf("a:%d:{", len(data))
	for key, value := range data {
		ret = ret + fmt.Sprintf("s:%d:\"%s\";", len(key), key)
		if valuemap, ok := value.(map[string]interface{}); ok {
			ret = ret + SerializePhp(valuemap)
		} else {
			valuestr := value.(string)
			ret = ret + fmt.Sprintf("s:%d:\"%s\";", len(valuestr), valuestr)
		}
	}
	ret = ret + "}"
	return ret
}

//channel超时时间
func Timeout(timeout time.Duration, ch chan bool) {
	time.Sleep(timeout)
	ch <- true
}

//判断是否在slice中
func InSlice(target interface{}, slice interface{}) (index int, err error) {
	t := reflect.ValueOf(slice).Kind().String()
	if t == "slice" {
		for i, val := range t {
			if target == val {
				index = i
				return
			}
		}

		err = errors.New("target not exists")
		return
	} else {
		err = errors.New("only support array now.")
		return
	}
}

/**
 * @des 生成唯一id(32位)
 *
 */
func GenUuid() string {
	nano := time.Now().UnixNano()
	rand.Seed(nano)
	rndNum := rand.Int63()
	hostName, _ := os.Hostname()
	return Md5(hostName + strconv.FormatInt(nano, 10) + strconv.FormatInt(rndNum, 10))
}

func Md5(text string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, text)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}

// RemoteIp 返回远程客户端的 IP，如 192.168.1.1
func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get(XRealIP); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get(XForwardedFor); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

// Ip2long 将 IPv4 字符串形式转为 uint32
func Ip2long(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

func Long2ip(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip>>24, ip<<8>>24, ip<<16>>24, ip<<24>>24)
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func IsPublicIP(req *http.Request) bool {
	ipstring := RemoteIp(req)
	IP := net.ParseIP(ipstring)
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

