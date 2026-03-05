package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type AccountInfo struct {
	User     string
	Cookie   string
	GameList []string
}

type CheckInResult struct {
	User     string
	GameCode string
	Success  bool
	Message  string
}

type ApiData struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
}

var apiMap = map[string]string{
	"gi":    "https://sg-hk4e-api.hoyolab.com/event/sol/sign?lang=zh-cn&act_id=e202102251931481",
	"hk3":   "https://sg-public-api.hoyolab.com/event/mani/sign?lang=zh-cn&act_id=e202110291205111",
	"hkrpg": "https://sg-public-api.hoyolab.com/event/luna/hkrpg/os/sign?lang=zh-cn&act_id=e202303301540311",
	"nxx":   "https://sg-public-api.hoyolab.com/event/luna/nxx/os/sign?lang=zh-cn&act_id=e202308141137581",
	"zzz":   "https://sg-public-api.hoyolab.com/event/luna/zzz/os/sign?lang=zh-cn&act_id=e202406031448091",
}

var gameName = map[string]string{
	"gi":    "Genshin",
	"hk3":   "Honkai_3",
	"hkrpg": "Star_Rail",
	"nxx":   "Tears_of_Themis",
	"zzz":   "Zenless_Zone_Zero",
}

var client = &http.Client{Timeout: 15 * time.Second}

func validateCookie(cookie string) bool {
	cookie = strings.TrimSpace(cookie)
	ltoken := strings.Contains(cookie, "ltoken_v2")
	ltuid := strings.Contains(cookie, "ltuid_v2")
	return ltoken && ltuid
}

func checkIn(userName, gameCode, cookie, url string) CheckInResult {
	result := CheckInResult{
		User:     userName,
		GameCode: gameCode,
	}
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("x-rpc-signgame", gameCode)
	req.Header.Set("x-rpc-client_type", "5")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://act.hoyolab.com")
	req.Header.Set("Origin", "https://act.hoyolab.com")

	resp, err := client.Do(req)
	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("网络请求失败: %v", err)
		return result
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data ApiData
	if err := json.Unmarshal(body, &data); err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("解析 JSON 失败: %v", err)
		return result
	}

	result.Success = (data.Retcode == 0)
	result.Message = data.Message
	return result
}

func logPrint(result CheckInResult) {
	icon := "✔️"
	if !result.Success {
		icon = "❌"
	}
	fmt.Printf("%s [%s] %s: %s\n", icon, result.User, gameName[result.GameCode], result.Message)
}

func main() {
	userCookie := strings.TrimSpace(os.Getenv("USER_COOKIE"))
	enabledGames := strings.TrimSpace(os.Getenv("ENABLED_GAMES"))
	if userCookie == "" {
		log.Fatal("环境变量 USER_COOKIE 未设置！")
	}
	if enabledGames == "" {
		log.Fatal("环境变量 ENABLED_GAMES 未设置！")
	}

	var cookieMap map[string]string
	var gamesMap map[string][]string
	if err := json.Unmarshal([]byte(userCookie), &cookieMap); err != nil {
		log.Fatal("环境变量 USER_COOKIE 不符合格式，请参考readme！", err)
	}
	if err := json.Unmarshal([]byte(enabledGames), &gamesMap); err != nil {
		log.Fatal("环境变量 ENABLED_GAMES 不符合格式，请参考readme！", err)
	}

	var accounts []AccountInfo
	for user, cookie := range cookieMap {
		if games, ok := gamesMap[user]; ok && len(games) > 0 {
			accounts = append(accounts, AccountInfo{
				User:     user,
				Cookie:   cookie,
				GameList: games,
			})
		}
	}

	var wg sync.WaitGroup
	fmt.Printf("🚀 开始为 %d 个账户执行自动签到任务...\n", len(accounts))
	for _, account := range accounts {
		wg.Add(1)
		go func(a AccountInfo) {
			defer wg.Done()
			if !validateCookie(a.Cookie) {
				fmt.Printf("%s的 Cookie 格式不正确，缺少 ltoken 或 ltuid...\n", a.User)
				return
			}

			for _, gameCode := range a.GameList {
				url, ok := apiMap[gameCode]
				if !ok {
					logPrint(CheckInResult{User: a.User, GameCode: gameCode, Success: false, Message: "游戏代码不合法"})
					continue
				}

				logPrint(checkIn(a.User, gameCode, a.Cookie, url))
				time.Sleep(time.Second)
			}
		}(account)
	}
	wg.Wait()

	fmt.Println("✨ 签到任务完成 ~")
}
