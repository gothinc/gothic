package httpclient

import "net/url"

/**
 * @desc Util
 * @author zhaojiangwei
 * @date 2018-12-25 18:03
 */

func encodeUrl(uri string) string {
	urlObj, err := url.Parse(uri)
	if err != nil{
		return ""
	}

	queryObj := urlObj.Query()
	urlObj.RawQuery = queryObj.Encode()
	return urlObj.String()
}