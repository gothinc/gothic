package library

import (
	"github.com/gothinc/gothic"
	"extension/qconf"
	"fmt"
	"strings"
)

/**
 * @desc Qconf.go
 * @author zhaojiangwei
 * @date 2018-12-27 17:38
 */

func GetConf(key string) map[string]string {
	idc := gothic.Config.GetString("application.cur_idc")

	value, err_conf := qconf.GetConf(key, idc)
	if err_conf != nil {
		err_msg := fmt.Sprintf("get qconf fail, key[%s], errmsg[%s]", key, err_conf)
		panic(err_msg)
	}
	ret := map[string]string{}

	vars := strings.Split(value, "|")
	for _, v := range vars {
		val := strings.Split(v, "=")
		ret[val[0]] = val[1]
	}

	return ret
}