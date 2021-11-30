package main

import (
	"encoding/base64"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"proxy_checker/configs"
	"proxy_checker/types"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func GetProxy() types.ProxyList {
	var proxy_list types.ProxyList
	jsonFile, err := ioutil.ReadFile("proxy.json")
	if err != nil {
		return proxy_list
	}

	err = json.Unmarshal([]byte(string(jsonFile)), &proxy_list)

	return proxy_list
}

func SetProxy(proxy types.ProxyList) {
	marshal, err := json.MarshalIndent(proxy, "", " ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile("proxy.json", marshal, 0777)
	if err != nil {
		return
	}
}


func main() {
	conf := configs.New()
	fmt.Println(conf)
	proxs := GetProxy()
	for i := 0; i < len(proxs.ProxyList); i++ {
		proxyStr := "http://" + proxs.ProxyList[i].Credentials + "@" + proxs.ProxyList[i].Address + ":" + proxs.ProxyList[i].Port
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			log.Println(err)
		}

		//creating the URL to be loaded through the proxy
		urlStr := conf.SiteAddress
		url, err := url.Parse(urlStr)
		if err != nil {
			log.Println(err)
		}

		//adding the proxy settings to the Transport object
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}

		//adding the Transport object to the http Client
		client := &http.Client{
			Transport: transport,
			Timeout: 5 * time.Second,
		}

		//generating the HTTP GET request
		request, err := http.NewRequest("GET", url.String(), nil)
		if err != nil {
			log.Println(err)
		}

		//adding proxy authentication
		auth := proxs.ProxyList[i].Credentials
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		request.Header.Add("Proxy-Authorization", basicAuth)
		//calling the URL
		response, err := client.Do(request)
		fmt.Println(proxs.ProxyList[i].Address)
		if err != nil {
			fmt.Println("Don't work or unpayed")
			proxs.ProxyList[i].Status = "Don't work or unpayed"
		} else if response.StatusCode == 403 {
			fmt.Println("Catch 403. Get ban from " + conf.SiteAddress)
			proxs.ProxyList[i].Status = "Catch 403. Get ban from " + conf.SiteAddress
		} else {
			fmt.Println("Work")
			proxs.ProxyList[i].Status = "Work"
		}
	}

	SetProxy(proxs)
	os.Exit(3)
}
