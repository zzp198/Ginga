package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

func main() {
	// 创建一个支持Cookie的Client
	cookieJar, _ := cookiejar.New(nil)
	cookie := &http.Cookie{
		Name:  "PHPSESSID",
		Value: "ibgfjvutgcor7lsnto2jenfh2s",
	}
	url, _ := url.Parse("https://free.vps.vc/vps-info")
	cookieJar.SetCookies(url, []*http.Cookie{cookie})

	client := &http.Client{
		Jar: cookieJar,
	}

	// 创建请求
	req, _ := http.NewRequest(http.MethodGet, "https://free.vps.vc/vps-info", nil)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	html, _ := io.ReadAll(resp.Body)

	regex, _ := regexp.Compile("Valid until</label>\\s*<div .*>\\s*(.*?)\\s*</div>")
	matches := regex.FindAllStringSubmatch(string(html), -1)
	for _, match := range matches {
		fmt.Println(match[1])
	}
}
