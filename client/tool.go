package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

type TmpMap struct {
	Keys    []int          `json:"keys"`
	Nodemap map[int]string `json:"map"`
}

// json转变为tmpMap
func JsonToTmpMap(jsonStr string) (*TmpMap, error) {
	m := TmpMap{}
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		fmt.Printf("[LOG_JsonToMap] Unmarshal with error: %+v\n", err)
		return nil, err
	}
	return &m, nil
}

// 解析获得value
func UnmarshalToValue(resp *http.Response) (result string) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[LOG_GET] read body failed")
		return
	}
	msg := make(map[string]string)
	err = json.Unmarshal(body, &msg)
	if err != nil {
		log.Println("[LOG_GET] Unmarshal failed", err.Error())
		return
	}
	return msg["msg"]
}

// 获取最接近的节点的地址
func (m *ClientMap) getUrl(key string) string {
	m.RLock()
	defer m.RUnlock()
	if len(m.Keys) == 0 { // 不符合条件
		return ""
	}
	hash := int(curHash([]byte(key)))
	// 二分查找最近的节点
	idx := sort.Search(len(m.Keys), func(i int) bool {
		return m.Keys[i] >= hash
	})
	idx = idx % len(m.Keys)
	// 返回地址
	return m.NodeMap[m.Keys[idx]]
}
