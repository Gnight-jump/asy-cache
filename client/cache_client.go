/*
	客户端示例：
		1. 配置服务中心
		2. 定时拉取服务中心的节点图
		3. 根据节点图查询，无法连接，则顺延，同时主动拉取新的
*/
package client

import (
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Hash 函数类型
type Hash func(data []byte) uint32

// 默认使用crc32算法
var (
	curHash    Hash   = crc32.ChecksumIEEE
	CenterPath string = ""

	// 错误显示
	NoCache    = errors.New("No cache found")
	SetFailed  = errors.New("Failed to set key-value pairs")
	DelFailed  = errors.New("Failed to del key-value pairs")
	ConnFailed = errors.New("Connect cache failed")
)

// ClientMap 使用map存储所有节点信息
type ClientMap struct {
	Keys         []int          // 存储节点信息（map的key）
	NodeMap      map[int]string // 存储整体哈希值
	sync.RWMutex                // 读写锁
}

// 新建hash
func New() *ClientMap {
	m := &ClientMap{
		Keys:    make([]int, 0),
		NodeMap: make(map[int]string),
	}
	// 定时任务
	m.FreshMap()
	go m.Timework()
	return m
}

// 定时任务，刷新map
func (m *ClientMap) Timework() {
	// 刷新时间:1min
	hourTicker := time.NewTicker(time.Minute * 1)
	go func() {
		for range hourTicker.C {
			m.FreshMap()
		}
	}()
}

// 刷新一致性哈希图 "/getMap"
func (m *ClientMap) FreshMap() {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("[LOG_FreshMap] recover error:", err)
		}
	}()

	resp, err := http.PostForm(CenterPath+"/getMap", url.Values{})
	if err != nil {
		log.Println("[LOG_FreshMap] connect", CenterPath+"/getMap", "failed")
		m.FreshMap()
	}

	// 解析文本
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[LOG_FreshMap] ReadAll error", err.Error())
	}
	jsonRes, err := JsonToTmpMap(string(body))
	if err != nil {
		fmt.Printf("[LOG_FreshMap] Convert json to map failed with error: %+v\n", err)
	}

	// 解析
	m.Keys = jsonRes.Keys
	m.NodeMap = jsonRes.Nodemap
}

// 获取key下的value "/get"
func (m *ClientMap) Get(key string) (result string) {
	v := ""
	// 寻找对应url
	if v = m.getUrl(key); v == "" {
		return
	}

	// 发起请求
	resp, err := http.PostForm(v+"/get", url.Values{
		"key": {key},
	})
	if err != nil {
		log.Println("[LOG_GET] connect", v, "failed")
		m.FreshMap()
		return
	}
	if resp.StatusCode != 200 {
		log.Println("[LOG_GET] get key no found")
		return
	}

	// 解析报文
	return UnmarshalToValue(resp)
}

// 设置键值对 "/set"
func (m *ClientMap) Set(key string, value string) error {
	v := ""
	// 寻找对应url
	if v = m.getUrl(key); v == "" {
		return NoCache
	}

	// 发起请求
	resp, err := http.PostForm(v+"/set", url.Values{
		"key":   {key},
		"value": {value},
	})
	if err != nil {
		log.Println("[LOG_SET] connect", v, "failed")
		m.FreshMap()
		return ConnFailed
	}
	if resp.StatusCode != 200 {
		log.Println("[LOG_SET] set key failed")
		return SetFailed
	}

	// 成功设置
	return nil
}

// 删除键值对 "/del"
func (m *ClientMap) Del(key string) error {
	v := ""
	// 寻找对应url
	if v = m.getUrl(key); v == "" {
		return NoCache
	}

	// 发起请求
	resp, err := http.PostForm(v+"/del", url.Values{
		"key": {key},
	})
	if err != nil {
		log.Println("[LOG_DEL] connect", v, "failed")
		m.FreshMap()
		return ConnFailed
	}
	if resp.StatusCode != 200 {
		log.Println("[LOG_DEL] del key failed")
		return DelFailed
	}

	// 成功设置
	return nil
}
