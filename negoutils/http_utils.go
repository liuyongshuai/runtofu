// @author      Liu Yongshuai<liuyongshuai@hotmail.com>
// @date        2018-11-22 18:37

package negoutils

import (
	"net/url"
	"strings"
)

// 对原始cookie进行五马分尸
func SplitRawCookie(ck string) (ret map[string]string) {
	ret = make(map[string]string)
	ck = strings.TrimSpace(ck)
	if len(ck) == 0 {
		return
	}
	kvs := strings.Split(ck, ";")
	if len(kvs) == 0 {
		return
	}
	for _, val := range kvs {
		val = strings.TrimSpace(val)
		if !strings.Contains(val, "=") {
			continue
		}
		ind := strings.Index(val, "=")
		k := strings.TrimSpace(val[0:ind])
		v := strings.TrimSpace(val[ind+1:])
		if len(k) == 0 || len(v) == 0 {
			continue
		}
		ret[k] = v
	}
	return
}

// 合并cookie
func JoinRawCookie(ck map[string]string) (ret string) {
	if ck == nil {
		return ""
	}
	if len(ck) == 0 {
		return ""
	}
	var tmp []string
	for k, v := range ck {
		tmp = append(tmp, k+"="+v)
	}
	ret = strings.Join(tmp, "; ")
	return
}

// urlencode()
func UrlEncode(str string) string {
	return url.QueryEscape(str)
}

// urldecode()
func UrlDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}

// 高仿PHP的rawurlencode()函数，在调M端的接口时有参数要求这么处理
func RawUrlEncode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// rawurldecode()
func RawUrlDecode(str string) (string, error) {
	return url.QueryUnescape(strings.Replace(str, "%20", "+", -1))
}

// parse_url()
// Parse a URL and return its components
// -1: all; 1: scheme; 2: host; 4: port; 8: user; 16: pass; 32: path; 64: query; 128: fragment
func ParseUrl(str string, component int) (map[string]string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return nil, err
	}
	if component == -1 {
		component = 1 | 2 | 4 | 8 | 16 | 32 | 64 | 128
	}
	var components = make(map[string]string)
	if (component & 1) == 1 {
		components["scheme"] = u.Scheme
	}
	if (component & 2) == 2 {
		components["host"] = u.Hostname()
	}
	if (component & 4) == 4 {
		components["port"] = u.Port()
	}
	if (component & 8) == 8 {
		components["user"] = u.User.Username()
	}
	if (component & 16) == 16 {
		components["pass"], _ = u.User.Password()
	}
	if (component & 32) == 32 {
		components["path"] = u.Path
	}
	if (component & 64) == 64 {
		components["query"] = u.RawQuery
	}
	if (component & 128) == 128 {
		components["fragment"] = u.Fragment
	}
	return components, nil
}

// http_build_query()
func HttpBuildQuery(queryData url.Values) string {
	return queryData.Encode()
}
