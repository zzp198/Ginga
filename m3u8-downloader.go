package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/yapingcat/gomedia/go-codec"
	"github.com/yapingcat/gomedia/go-mp4"
	"github.com/yapingcat/gomedia/go-mpeg2"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//func main() {
//
//	download("https://vip.lz-cdn5.com/20220328/1825_7132edf8/1200k/hls/index.m3u8", "ylxq.mp4")
//}

func download(rawUrl, name string) (bool, error) {
	dir := rawUrl[:strings.LastIndex(rawUrl, "/")]

	fmt.Println(dir)

	c := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := c.Get(rawUrl)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	m3u8 := string(data)

	re, _ := regexp.Compile("#EXTINF:([\\d.]+)([,\n]+)(?P<ru>.+)")

	play_list := make([]string, 0)

	for _, v := range re.FindAllString(m3u8, -1) {
		ru := re.FindStringSubmatch(v)[re.SubexpIndex("ru")]
		play, _ := url.JoinPath(dir, ru)
		play_list = append(play_list, play)
	}

	mp4f, _ := os.Create("temp/" + name)
	defer mp4f.Close()

	muxer, _ := mp4.CreateMp4Muxer(mp4f)
	vtid := muxer.AddVideoTrack(mp4.MP4_CODEC_H264)
	atid := muxer.AddAudioTrack(mp4.MP4_CODEC_AAC)

	demuxer := mpeg2.NewTSDemuxer()
	var OnFrameErr error
	var audioTimestamp uint64 = 0
	aacSampleRate := -1
	demuxer.OnFrame = func(cid mpeg2.TS_STREAM_TYPE, frame []byte, pts uint64, dts uint64) {
		if OnFrameErr != nil {
			return
		}
		if cid == mpeg2.TS_STREAM_AAC {
			audioTimestamp = pts
			codec.SplitAACFrame(frame, func(aac []byte) {
				if aacSampleRate == -1 {
					adts := codec.NewAdtsFrameHeader()
					adts.Decode(aac)
					aacSampleRate = codec.AACSampleIdxToSample(int(adts.Fix_Header.Sampling_frequency_index))
				}
				err = muxer.Write(atid, aac, audioTimestamp, audioTimestamp)
				audioTimestamp += uint64(1024 * 1000 / aacSampleRate) //每帧aac采样固定为1024。aac_sampleRate 为采样率
				if err != nil {
					OnFrameErr = err
					return
				}
			})
		} else if cid == mpeg2.TS_STREAM_H264 {
			err = muxer.Write(vtid, frame, pts, dts)
			if err != nil {
				OnFrameErr = err
				return
			}
		} else {
			OnFrameErr = errors.New("unknown cid " + strconv.Itoa(int(cid)))
			return
		}
	}

	for i := 0; i < len(play_list); i++ {
		play := play_list[i]

		fmt.Println(play)

		playresp, _ := http.Get(play)

		rd, _ := io.ReadAll(playresp.Body)

		err = demuxer.Input(bytes.NewReader(rd))

		playresp.Body.Close()
	}

	err = muxer.WriteTrailer()

	err = mp4f.Sync()

	return true, nil
}

//// 该技术基本原理是将视频文件或视频流切分成小片(ts)并建立索引文件(m3u8)。
//// 支持的视频流编码为H.264，音频流编码为AAC。
//func download(ru, name string) (bool, error) {
//
//	u, _ := url.Parse(ru)
//	host := u.Scheme + "://" + u.Host + u.EscapedPath()
//
//	play_list := make(map[string]string)
//
//	c := http.Client{
//		Transport: &http.Transport{
//			TLSClientConfig: &tls.Config{
//				InsecureSkipVerify: true,
//			},
//		},
//	}
//
//	resp, err := c.Get(ru)
//	if err != nil {
//		return false, err
//	}
//	defer resp.Body.Close()
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return false, err
//	}
//
//	pwd, _ := os.Getwd()
//	download_dir := filepath.Join(pwd, name)
//
//	_, err = os.Stat(download_dir)
//	if err != nil && os.IsNotExist(err) {
//		_ = os.MkdirAll(download_dir, os.ModePerm)
//	}
//
//	return true, nil
//}
//
//// 获取m3u8加密的密钥
//func getM3u8Key(host, html string) (key string) {
//	lines := strings.Split(html, "\n")
//	key = ""
//	for _, line := range lines {
//		if strings.Contains(line, "#EXT-X-KEY") {
//			uri_pos := strings.Index(line, "URI")
//			quotation_mark_pos := strings.LastIndex(line, "\"")
//			key_url := strings.Split(line[uri_pos:quotation_mark_pos], "\"")[1]
//			if !strings.Contains(line, "http") {
//				key_url = fmt.Sprintf("%s/%s", host, key_url)
//			}
//			res, err := grequests.Get(key_url, ro)
//			checkErr(err)
//			if res.StatusCode == 200 {
//				key = res.String()
//			}
//		}
//	}
//	return
//}
//
//func getTsList(host, body string) (tsList []TsInfo) {
//	lines := strings.Split(body, "\n")
//	index := 0
//	var ts TsInfo
//	for _, line := range lines {
//		if !strings.HasPrefix(line, "#") && line != "" {
//			//有可能出现的二级嵌套格式的m3u8,请自行转换！
//			index++
//			if strings.HasPrefix(line, "http") {
//				ts = TsInfo{
//					Name: fmt.Sprintf(TS_NAME_TEMPLATE, index),
//					Url:  line,
//				}
//				tsList = append(tsList, ts)
//			} else {
//				ts = TsInfo{
//					Name: fmt.Sprintf(TS_NAME_TEMPLATE, index),
//					Url:  fmt.Sprintf("%s/%s", host, line),
//				}
//				tsList = append(tsList, ts)
//			}
//		}
//	}
//	return
//}
//
//// 下载ts文件
//// @modify: 2020-08-13 修复ts格式SyncByte合并不能播放问题
//func downloadTsFile(ts TsInfo, download_dir, key string, retries int) {
//	defer func() {
//		if r := recover(); r != nil {
//			//fmt.Println("网络不稳定，正在进行断点持续下载")
//			downloadTsFile(ts, download_dir, key, retries-1)
//		}
//	}()
//	curr_path_file := fmt.Sprintf("%s/%s", download_dir, ts.Name)
//	if isExist, _ := pathExists(curr_path_file); isExist {
//		//logger.Println("[warn] File: " + ts.Name + "already exist")
//		return
//	}
//	res, err := grequests.Get(ts.Url, ro)
//	if err != nil || !res.Ok {
//		if retries > 0 {
//			downloadTsFile(ts, download_dir, key, retries-1)
//			return
//		} else {
//			//logger.Printf("[warn] File :%s", ts.Url)
//			return
//		}
//	}
//	// 校验长度是否合法
//	var origData []byte
//	origData = res.Bytes()
//	contentLen := 0
//	contentLenStr := res.Header.Get("Content-Length")
//	if contentLenStr != "" {
//		contentLen, _ = strconv.Atoi(contentLenStr)
//	}
//	if len(origData) == 0 || (contentLen > 0 && len(origData) < contentLen) || res.Error != nil {
//		//logger.Println("[warn] File: " + ts.Name + "res origData invalid or err：", res.Error)
//		downloadTsFile(ts, download_dir, key, retries-1)
//		return
//	}
//	// 解密出视频 ts 源文件
//	if key != "" {
//		//解密 ts 文件，算法：aes 128 cbc pack5
//		origData, err = AesDecrypt(origData, []byte(key))
//		if err != nil {
//			downloadTsFile(ts, download_dir, key, retries-1)
//			return
//		}
//	}
//	// https://en.wikipedia.org/wiki/MPEG_transport_stream
//	// Some TS files do not start with SyncByte 0x47, they can not be played after merging,
//	// Need to remove the bytes before the SyncByte 0x47(71).
//	syncByte := uint8(71) //0x47
//	bLen := len(origData)
//	for j := 0; j < bLen; j++ {
//		if origData[j] == syncByte {
//			origData = origData[j:]
//			break
//		}
//	}
//	ioutil.WriteFile(curr_path_file, origData, 0666)
//}
//
//// downloader m3u8 下载器
//func downloader(tsList []TsInfo, maxGoroutines int, downloadDir string, key string) {
//	retry := 5 //单个 ts 下载重试次数
//	var wg sync.WaitGroup
//	limiter := make(chan struct{}, maxGoroutines) //chan struct 内存占用 0 bool 占用 1
//	tsLen := len(tsList)
//	downloadCount := 0
//	for _, ts := range tsList {
//		wg.Add(1)
//		limiter <- struct{}{}
//		go func(ts TsInfo, downloadDir, key string, retryies int) {
//			defer func() {
//				wg.Done()
//				<-limiter
//			}()
//			downloadTsFile(ts, downloadDir, key, retryies)
//			downloadCount++
//			DrawProgressBar("Downloading", float32(downloadCount)/float32(tsLen), PROGRESS_WIDTH, ts.Name)
//			return
//		}(ts, downloadDir, key, retry)
//	}
//	wg.Wait()
//}
//
//func checkTsDownDir(dir string) bool {
//	if isExist, _ := pathExists(filepath.Join(dir, fmt.Sprintf(TS_NAME_TEMPLATE, 0))); !isExist {
//		return true
//	}
//	return false
//}
//
//// 合并ts文件
//func mergeTs(downloadDir string) string {
//	mvName := downloadDir + ".mp4"
//	outMv, _ := os.Create(mvName)
//	defer outMv.Close()
//	writer := bufio.NewWriter(outMv)
//	err := filepath.Walk(downloadDir, func(path string, f os.FileInfo, err error) error {
//		if f == nil {
//			return err
//		}
//		if f.IsDir() || filepath.Ext(path) != ".ts" {
//			return nil
//		}
//		bytes, _ := ioutil.ReadFile(path)
//		_, err = writer.Write(bytes)
//		return err
//	})
//	checkErr(err)
//	_ = writer.Flush()
//	os.RemoveAll(downloadDir)
//	return mvName
//}
//
//// ============================== 加解密相关 ==============================
//
//func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
//	padding := blockSize - len(ciphertext)%blockSize
//	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
//	return append(ciphertext, padtext...)
//}
//
//func PKCS7UnPadding(origData []byte) []byte {
//	length := len(origData)
//	unpadding := int(origData[length-1])
//	return origData[:(length - unpadding)]
//}
//
//func AesEncrypt(origData, key []byte, ivs ...[]byte) ([]byte, error) {
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		return nil, err
//	}
//	blockSize := block.BlockSize()
//	var iv []byte
//	if len(ivs) == 0 {
//		iv = key
//	} else {
//		iv = ivs[0]
//	}
//	origData = PKCS7Padding(origData, blockSize)
//	blockMode := cipher.NewCBCEncrypter(block, iv[:blockSize])
//	crypted := make([]byte, len(origData))
//	blockMode.CryptBlocks(crypted, origData)
//	return crypted, nil
//}
//
//// #EXTINF:([\d\.]+),(?P<url>.*)
//func AesDecrypt(crypted, key []byte, ivs ...[]byte) ([]byte, error) {
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		return nil, err
//	}
//	blockSize := block.BlockSize()
//	var iv []byte
//	if len(ivs) == 0 {
//		iv = key
//	} else {
//		iv = ivs[0]
//	}
//	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
//	origData := make([]byte, len(crypted))
//	blockMode.CryptBlocks(origData, crypted)
//	origData = PKCS7UnPadding(origData)
//	return origData, nil
//}
