package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http" 
	"sort"
	"strconv"
	"sync"
	"time"
	"bytes"

	"github.com/deepch/vdk/av"
)

var (
	Success                         = "success"
	ErrorStreamNotFound             = errors.New("stream not found")
	ErrorStreamAlreadyExists        = errors.New("stream already exists")
	ErrorStreamChannelAlreadyExists = errors.New("stream channel already exists")
	ErrorStreamNotHLSSegments       = errors.New("stream hls not ts seq found")
	ErrorStreamNoVideo              = errors.New("stream no video")
	ErrorStreamNoClients            = errors.New("stream no clients")
	ErrorStreamRestart              = errors.New("stream restart")
	ErrorStreamStopCoreSignal       = errors.New("stream stop core signal")
	ErrorStreamStopRTSPSignal       = errors.New("stream stop rtsp signal")
	ErrorStreamChannelNotFound      = errors.New("stream channel not found")
	ErrorStreamChannelCodecNotFound = errors.New("stream channel codec not ready, possible stream offline")
	ErrorStreamsLen0                = errors.New("streams len zero")
)

//Config global
var Config = loadConfig()

//ConfigST struct
type ConfigST struct {
	mutex   sync.RWMutex
	Server  ServerST            `json:"server"`
	Streams map[string]StreamST `json:"streams"`
}

//ServerST struct
type ServerST struct {
	HTTPPort string `json:"http_port"`
}

//StreamST struct
type StreamST struct {
	URL              string          `json:"url"`
	Status           bool            `json:"status"`
	OnDemand         bool            `json:"on_demand"`
	RunLock          bool            `json:"-"`
	hlsSegmentNumber int             `json:"-"`
	hlsSegmentBuffer map[int]Segment `json:"-"`
	Codecs           []av.CodecData
	Cl               map[string]viewer
}

//Segment HLS cache section
type Segment struct {
	dur  time.Duration
	data []*av.Packet
}

type viewer struct {
	c chan av.Packet
}

func (element *ConfigST) RunIFNotRun(uuid string) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	if tmp, ok := element.Streams[uuid]; ok {
		if tmp.OnDemand && !tmp.RunLock {
			tmp.RunLock = true
			element.Streams[uuid] = tmp
			go RTSPWorkerLoop(uuid, tmp.URL, tmp.OnDemand)
		}
	}
}

func (element *ConfigST) RunUnlock(uuid string) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	if tmp, ok := element.Streams[uuid]; ok {
		if tmp.OnDemand && tmp.RunLock {
			tmp.RunLock = false
			element.Streams[uuid] = tmp
		}
	}
}

func (element *ConfigST) HasViewer(uuid string) bool {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	if tmp, ok := element.Streams[uuid]; ok && len(tmp.Cl) > 0 {
		return true
	}
	return false
}

// //声明一个名为 loadConfig 的函数，返回值是 *ConfigST 类型的指针（即返回一个配置对象的地址）
// func loadConfig() *ConfigST {
// 	//创建一个 ConfigST 类型的变量 tmp，用来暂存加载出来的配置内容。
// 	var tmp ConfigST
// 	//读取当前目录下的 config.json 文件内容，结果保存在 data 中，若失败会返回 err。
// 	data, err := ioutil.ReadFile("config.json")
// 	//如果读取文件时出错，就打印错误并终止程序（log.Fatalln 会输出日志并调用 os.Exit(1) 退出程序）。
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	//将读取到的 JSON 数据解析（反序列化）到 tmp 对象中（data → tmp），相当于把 JSON 数据转为结构体。
// 	err = json.Unmarshal(data, &tmp)
// 	//如果解析出错，也会打印错误并退出程序。
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	//遍历解析得到的所有视频流:
// 	for i, v := range tmp.Streams {
// 		//为每个流的 Cl（客户端列表）创建一个新的空 map。
// 		v.Cl = make(map[string]viewer)
// 		//为每个流的 hlsSegmentBuffer（HLS 缓存段）也创建一个新的空 map。
// 		v.hlsSegmentBuffer = make(map[int]Segment)
// 		//然后把修改后的流结构体重新赋值回 tmp.Streams[i]
// 		tmp.Streams[i] = v
// 	}
// 	return &tmp
// }

func loadConfig() *ConfigST {
	var tmp ConfigST
	tmp.Server.HTTPPort = ":8083"
	tmp.Streams = make(map[string]StreamST)

	// 从后端动态获取 cameraCodes 列表
	resp, err := http.Get("http://localhost:3000/cameras") // 替换成 https://10.70.37.12:18531/device/deviceList/v1.0?deviceType=2&fromIndex=1&toIndex=2000
	if err != nil {
		log.Printf("获取 cameraCode 列表失败: %v", err)
		return &tmp
	}
	defer resp.Body.Close()

	var cameraCodes []string
	if err := json.NewDecoder(resp.Body).Decode(&cameraCodes); err != nil {
		log.Printf("解析 cameraCode 列表失败: %v", err)
		return &tmp
	}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, code := range cameraCodes {
		reqBody := map[string]string{"cameraCode": code}
		jsonBody, _ := json.Marshal(reqBody)

		req, err := http.NewRequest("POST", "http://localhost:3000/video", bytes.NewBuffer(jsonBody)) // https://10.70.37.12:18531/video/rtspurl/v1.0
		if err != nil {
			log.Printf("构建请求失败（%s）: %v", code, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("请求失败（%s）: %v", code, err)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("读取响应失败（%s）: %v", code, err)
			continue
		}

		var result struct {
			RTSPURL string `json:"rtspURL"`
		}
		if err := json.Unmarshal(body, &result); err != nil || result.RTSPURL == "" {
			log.Printf("解析 JSON 或无效返回（%s）: %v", code, err)
			continue
		}

		tmp.Streams[code] = StreamST{
			URL:              result.RTSPURL,
			OnDemand:         true, // false 表示“非按需”，服务器启动后立即开始拉流 //true 表示“按需”，只有在有用户访问播放地址时才启动拉流
			Cl:               make(map[string]viewer),
			hlsSegmentBuffer: make(map[int]Segment),
		}

        log.Printf("尝试请求 cameraCode=%s 的 RTSP 地址", code)
		log.Printf("已添加摄像头: %s -> %s", code, result.RTSPURL)
	}

	return &tmp
}

func (element *ConfigST) cast(uuid string, pck av.Packet) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	for _, v := range element.Streams[uuid].Cl {
		if len(v.c) < cap(v.c) {
			v.c <- pck
		}
	}
}

func (element *ConfigST) ext(suuid string) bool {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	_, ok := element.Streams[suuid]
	return ok
}

func (element *ConfigST) coAd(suuid string, codecs []av.CodecData) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	t := element.Streams[suuid]
	t.Codecs = codecs
	element.Streams[suuid] = t
}

func (element *ConfigST) coGe(suuid string) []av.CodecData {
	for i := 0; i < 100; i++ {
		element.mutex.RLock()
		tmp, ok := element.Streams[suuid]
		element.mutex.RUnlock()
		if !ok {
			return nil
		}
		if tmp.Codecs != nil {
			return tmp.Codecs
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}

func (element *ConfigST) clAd(suuid string) (string, chan av.Packet) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	cuuid := pseudoUUID()
	ch := make(chan av.Packet, 100)
	element.Streams[suuid].Cl[cuuid] = viewer{c: ch}
	return cuuid, ch
}

func (element *ConfigST) list() (string, []string) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	var res []string
	var fist string
	for k := range element.Streams {
		if fist == "" {
			fist = k
		}
		res = append(res, k)
	}
	return fist, res
}
func (element *ConfigST) clDe(suuid, cuuid string) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	delete(element.Streams[suuid].Cl, cuuid)
}

func pseudoUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}

//StreamHLSAdd add hls seq to buffer
func (obj *ConfigST) StreamHLSAdd(uuid string, val []*av.Packet, dur time.Duration) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		tmp.hlsSegmentNumber++
		tmp.hlsSegmentBuffer[tmp.hlsSegmentNumber] = Segment{data: val, dur: dur}
		if len(tmp.hlsSegmentBuffer) >= 6 {
			delete(tmp.hlsSegmentBuffer, tmp.hlsSegmentNumber-6-1)
		}
		obj.Streams[uuid] = tmp
	}
}

//StreamHLSm3u8 get hls m3u8 list
func (obj *ConfigST) StreamHLSm3u8(uuid string) (string, int, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		var out string
		//TODO fix  it
		out += "#EXTM3U\r\n#EXT-X-TARGETDURATION:4\r\n#EXT-X-VERSION:4\r\n#EXT-X-MEDIA-SEQUENCE:" + strconv.Itoa(tmp.hlsSegmentNumber) + "\r\n"
		var keys []int
		for k := range tmp.hlsSegmentBuffer {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		var count int
		for _, i := range keys {
			count++
			out += "#EXTINF:" + strconv.FormatFloat(tmp.hlsSegmentBuffer[i].dur.Seconds(), 'f', 1, 64) + ",\r\nsegment/" + strconv.Itoa(i) + "/file.ts\r\n"

		}
		return out, count, nil
	}
	return "", 0, ErrorStreamNotFound
}

//StreamHLSTS send hls segment buffer to clients
func (obj *ConfigST) StreamHLSTS(uuid string, seq int) ([]*av.Packet, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if buf, ok := tmp.hlsSegmentBuffer[seq]; ok {
			return buf.data, nil
		}
	}
	return nil, ErrorStreamNotFound
}

//StreamHLSFlush delete hls cache
func (obj *ConfigST) StreamHLSFlush(uuid string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		tmp.hlsSegmentBuffer = make(map[int]Segment)
		tmp.hlsSegmentNumber = 0
		obj.Streams[uuid] = tmp
	}
}

//stringToInt convert string to int if err to zero
func stringToInt(val string) int {
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return i
}
