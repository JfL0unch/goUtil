package utils

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func HttpGet(api string) ([]byte, error) {
	resp, err := http.Get(api)
	if err != nil {
		return nil, errors.Wrap(err, "Do")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.Errorf("themis error with code %d,%s", resp.StatusCode, body)
	}
	m, e := ioutil.ReadAll(resp.Body)

	return m, e
}

func HttpPost(url, contentType string, data []byte) ([]byte, error) {
	resp, err := http.Post(url, contentType, bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "Do")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.Errorf("error with code %d,%s", resp.StatusCode, body)
	}
	m, e := ioutil.ReadAll(resp.Body)

	return m, e
}

type KV struct {
	k string
	v string
}
type Params map[string]string

func (p Params) Add(key, val string) {
	p[key] = val
}

// QueryStr
// escaped 是否转义字符
// sorted 是否按字典排序
func (p Params) QueryStr(sorted, escaped bool) string {
	var qs string
	res := p.KVs(sorted)

	for _, item := range res {
		k := item.k
		v := item.v

		if escaped {
			k = url.QueryEscape(k)
			v = url.QueryEscape(v)
		}

		if len(qs) == 0 {
			qs = fmt.Sprintf("%s=%s", k, v)
		} else {
			qs = fmt.Sprintf("%s&%s=%s", qs, k, v)
		}
	}
	return qs
}

func (p Params) Keys() []string {

	if p == nil {
		return nil
	}
	res := make([]string, 0, len(p))
	for k, _ := range p {
		if k != "" {
			res = append(res, k)
		}
	}
	return res
}

func (p Params) SortedKeys() []string {
	keys := make([]string, 0, len(p))
	for k, _ := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (p Params) KVs(sorted bool) []KV {
	if p == nil {
		return nil
	}

	res := make([]KV, 0)

	keys := make([]string, 0)
	if sorted {
		keys = p.SortedKeys()
	} else {
		keys = p.Keys()
	}

	for _, k := range keys {
		res = append(res, KV{
			k: k,
			v: p[k],
		})
	}
	return res
}

func (p Params) Vals(sorted, escaped bool) string {
	if p == nil {
		return ""
	}

	buf := strings.Builder{}

	keys := make([]string, 0)
	if sorted {
		keys = p.SortedKeys()
	} else {
		keys = p.Keys()
	}

	for _, k := range keys {
		v := p[k]
		if escaped {
			buf.WriteString(url.QueryEscape(v))
		} else {
			buf.WriteString(v)
		}
	}
	return buf.String()
}

type ParamItem struct {
	Key string
	Val []string
}

// deleteKeys 要删除的key
func toArray(params url.Values, deleteKeys []string) []ParamItem {

	res := make([]ParamItem, 0)

	for _, key := range deleteKeys {
		if _, ok := params[key]; ok {
			delete(params, key)
		}
	}

	for k, v := range params {
		res = append(res, ParamItem{
			Key: k,
			Val: v,
		})
	}
	return res
}

// SortUrlValues 根据 url query parameters 的 map 进行排序，排除掉 指定key 后生成字符串。
// deleteKeys 要删除的keys
func SortUrlValues(params url.Values, deleteKeys []string) string {
	lineList := make([]string, 0, len(params))
	ms := toArray(params, deleteKeys)
	sort.Slice(ms, func(i, j int) bool {
		if ms[i].Key < ms[j].Key {
			return true
		} else {
			return false
		}
	})
	for _, item := range ms {
		lineList = append(lineList, fmt.Sprintf("%s=%s", item.Key, item.Val[0]))
	}

	// build canonical query string
	canonicalString := strings.Join(lineList, "&")
	return canonicalString
}
