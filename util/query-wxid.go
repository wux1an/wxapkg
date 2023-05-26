package util

import (
	"bytes"
	"encoding/json"
	"errors"
	ua "github.com/wux1an/fake-useragent"
	"io"
	"net/http"
	"os"
)

var cachedWxid = make(map[string]WxidInfo)

const CachePath = "wxid.json"

func init() {
	loadWxidCache()
}

func loadWxidCache() {
	if file, err := os.ReadFile(CachePath); err != nil {
		return
	} else {
		if err := json.Unmarshal(file, &cachedWxid); err != nil {
			return
		}
	}
}

func saveCache() {
	data, err := json.MarshalIndent(cachedWxid, "", "  ")
	if err != nil {
		return
	}

	_ = os.WriteFile(CachePath, data, 0600)
}

var WxidQuery = &queryWxid{}

type WxidInfo struct {
	Wxid     string `json:"-"` // not marshal
	Location string `json:"-"` // not marshal
	Error    string `json:"-"` // not marshal

	Nickname      string `json:"nickname"`
	Username      string `json:"username"`
	Description   string `json:"description"`
	Avatar        string `json:"avatar"`
	UsesCount     string `json:"uses_count"`
	PrincipalName string `json:"principal_name"`
}

func (w WxidInfo) Json() string {
	data, _ := json.MarshalIndent(w, "", "  ")
	return string(data)
}

type queryWxid struct {
}

type queryWxidReqBody struct {
	Appid string `json:"appid"`
}

type queryWxidRespBody struct {
	Code   int      `json:"code"`
	Errors string   `json:"errors"`
	Data   WxidInfo `json:"data"`
}

func (*queryWxid) Query(wxid string) (WxidInfo, error) {
	info, exist := cachedWxid[wxid]
	if exist {
		return info, nil
	}

	body, _ := json.Marshal(queryWxidReqBody{Appid: wxid})
	req, err := http.NewRequest("POST", "https://kainy.cn/api/weapp/info/", bytes.NewReader(body))
	if err != nil {
		return WxidInfo{}, err
	}

	req.Header.Set("User-Agent", ua.Random())
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return WxidInfo{}, err
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return WxidInfo{}, err
	}

	var respObj queryWxidRespBody
	if err := json.Unmarshal(all, &respObj); err != nil {
		return WxidInfo{}, err
	}

	if respObj.Code != 0 {
		return WxidInfo{}, errors.New(respObj.Errors)
	}

	// cache this record
	cachedWxid[wxid] = respObj.Data
	saveCache()

	return respObj.Data, nil
}
