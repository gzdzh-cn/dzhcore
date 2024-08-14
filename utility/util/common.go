package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// const API_KEY = "b3URo8wVBaUTbQyZENVXrrUz"
// const SECRET_KEY = "f5cRL3G1oaRW8ds24ysPiluiLtliNed0"

/**
 * 使用 AK，SK 生成鉴权签名（Access Token）
 * @return string 鉴权签名信息（Access Token）
 */
func GetAccessToken(API_KEY string, SECRET_KEY string) string {

	url := "https://aip.baidubce.com/oauth/2.0/token"
	postData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", API_KEY, SECRET_KEY)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	accessTokenObj := map[string]string{}
	err = json.Unmarshal([]byte(body), &accessTokenObj)
	if err != nil {
		return ""
	}
	return accessTokenObj["access_token"]
}
