package postSample

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	"time"
)

var DEBUG_KBN = 0
var params = map[string]string{}

type BasicAuth struct{
	Id string
	Pass string
}

// とにかく最初に呼び出す
func Initialize() {
	params = map[string]string{}
}

func SetParams(name string, val string) {
	params[name] = val;
}

func HttpPost(api_url string) string {
	fmt.Println("request_url : " + getHttpUrl(api_url))

	// 送信用パラメーターを作成
	bparams := &bytes.Buffer{}
	writer := multipart.NewWriter(bparams)

	// パラメーターを書き込んでいく
	for key, val := range params {
		fmt.Println(key + ":" + val)
		file, err := os.Open(val)
		// エラーだった場合ファイルではない
		if err != nil {
			// 通常のパラメーターとして書き込み
			_ = writer.WriteField(key, val)
		} else {
			fmt.Println("is File")
			// ファイルとして書き込み
			part, err := writer.CreateFormFile(key, filepath.Base(val))
			if err != nil {
				return getErrorJSON(901, "ファイルの送信に失敗しました")
			}
			_, err = io.Copy(part, file)
		}
	}

	err := writer.Close()
	if err != nil {
		return getErrorJSON(902, "通信に失敗しました")
	}

	req, err := http.NewRequest(
        "POST",
        getHttpUrl(api_url),
        bparams)

	if err != nil {
		return getErrorJSON(903, "通信に失敗しました")
	}

	// 持続接続を設定
	req.Header.Set("Connection", "Keep-Alive")

	// Content-Type 設定
	req.Header.Add("Content-Type", writer.FormDataContentType())

    // Basic認証がある場合はコメントを外す
//    req.Header.Add("Authorization", getBasicAuth(auth_id, auth_pass))

    // タイムアウトを15秒に設定
    client := &http.Client{ Timeout: time.Duration(15 * time.Second) }
    resp, err := client.Do(req)
    if err != nil {
    	fmt.Println(err.Error())
    	return getErrorJSON(500, err.Error())
    }

    defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return getErrorJSON(904, err.Error())
	}

	// ステータスによるエラーハンドリング(正常終了じゃなければ)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("%d", resp.StatusCode)
		return getErrorJSON(resp.StatusCode, resp.Status)
    }

	return string(body)
}

// URLにプロトコルとドメインを追加する
func getHttpUrl(url string) string {
	// URLが空だった場合
	if len(url) <= 0 {
		return url;
	}

	HTTP := [2]string{"http://", "http://"}
	HOST_NAME := [2]string{"it-wiseman.com:8080", "dev.it-wiseman.com:8080"}

	// urlにhttp://もしくはhttps://が含まれていない場合
	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		return HTTP[DEBUG_KBN] + HOST_NAME[DEBUG_KBN] + url
	}

	return url
}

// Base64に変換し、Basic認証用のパスを作る
func getBasicAuth(key string, pass string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(url.QueryEscape(key) + ":" + url.QueryEscape(pass)))
}

// エラーデータの作成
func getErrorJSON(code int, message string) string {
	return "{\"header\":{\"code\":" + strconv.Itoa(code) + ", \"message\":\"" + message + "}, \"body\":{}}"
}
