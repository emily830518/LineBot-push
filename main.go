package main
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"strings"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/go-redis/redis"
	// "strconv"
	"time"
)

type device struct {
	// Gps_num float32 `json:"gps_num"`
	// App string `json:"app"`
	// Gps_alt float32 `json:"gps_alt"`
	// Fmt_opt int `json:"fmt_opt"`
	// Device string `json:"device"`
	// S_d2 float32 `json:"s_d2"`
	// S_d0 float32 `json:"s_d0"`
	// S_d1 float32 `json:"s_d1"`
	// S_h0 float32 `json:"s_h0"`
	// SiteName string `json:"SiteName"`
	// Gps_fix float32 `json:"gps_fix"`
	// Ver_app string `json:"ver_app"`
	// Gps_lat float32 `json:"gps_lat"`
	// S_t0 float32 `json:"s_t0"`
	// Timestamp string `json:"timestamp"`
	// Gps_lon float32 `json:"gps_lon"`
	// Date string `json:"date"`
	// Tick float32 `json:"tick"`
	Device_id string `json:"device_id"`
	Rate float32 `json:"rate"`
	// S_1 float32 `json:"s_1"`
	// S_0 float32 `json:"s_0"`
	// S_3 float32 `json:"s_3"`
	// S_2 float32 `json:"s_2"`
	// Ver_format string `json:"ver_format"`
	// Time string `json:"time"`
}

type airbox struct {
	Source string `json:"source"`
	Feeds []device `json:"feeds"`
	Version string `json:"version"`
}

var bot *linebot.Client
var indoor_json airbox
var	client=redis.NewClient(&redis.Options{
		Addr:"hipposerver.ddns.net:6379",
		Password:"",
		DB:0,
	})

func periodicFunc(tick time.Time){
	url := "https://data.lass-net.org/data/device_indoor.json"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	errs := json.Unmarshal(body, &indoor_json)
	if errs != nil {
		fmt.Println(errs)
	}
	var err error
	bot, err = linebot.New("e07874bf26b33e5397a0de0278557807", "yk6oBgGuxjuya/Ufrm53SSCotxktUAtA96rjRBlF3+G7muhj/iqAi8w+SeX8/IpaYlUQGIA1NdQ69xxuGrO/aKbmkZsaRAiETLwpSgCfkJ0NmiKuOSz67vHSOFNI+8EECWkvp1l7wvImemYa0SDwiwdB04t89/1O/w1cDnyilFU=")
	log.Println("Bot:", bot, " err:", err)
	// http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
	for i:=0; i<len(indoor_json.Feeds); i++ {
		val, err:=client.Get(indoor_json.Feeds[i].Device_id).Result()
		if err!=nil{
			continue
		}
		stringSlice:=strings.Split(val,",")
		for j:=0; j<len(stringSlice); j++ {
			_,_=bot.PushMessage(stringSlice[j], linebot.NewTextMessage("您訂閱的AirBox"+indoor_json.Feeds[i].Device_id+"被判定放置於室內，請將其移到通風良好的室外！謝謝！")).Do()
		}
	}
	// _,_=bot.PushMessage("U3617adbdd46283d7e859f36302f4f471", linebot.NewTextMessage("hi!")).Do()
    // fmt.Println("Tick at: ", tick)
}

func main() {
	ticker := time.NewTicker(1 * time.Minute)
	go func(){
        for t := range ticker.C {
            //Call the periodic function here.
            periodicFunc(t)
        }
    }()
    quit := make(chan bool, 1)
    <-quit
}	