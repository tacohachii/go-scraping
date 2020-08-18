package main

import (
		"fmt"
    "log"
    "time"
		"os"
		"bufio"
		"strconv"
		"strings"
		"net/http"
		"github.com/PuerkitoBio/goquery"
)

func ExtractUrlRegExp(str string) string {
		arr := strings.Split(str, "'")
		return arr[1]
}

func CapitalNumRegExp(str string) string {
		arr := strings.Split(strings.Split(str, "百万円")[0], ",")
		return strings.Join(arr, "")
}

func main() {
		rootUrl := "https://www.bizreach.jp"
		url := rootUrl + "/company"              // アクセスするURL

    file, err := os.Create(`bizreach.csv`)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    count := 0                              // いくつアクセスが成功したかをカウント
    start := time.Now()                     // 処理にかかった時間を測定
    for maxRequest := 0; maxRequest < 28; maxRequest ++ {     // リクエストを送る
				if maxRequest != 0 {
					page := strconv.Itoa(maxRequest + 1)
					url = rootUrl + "/company/c/?pageSize=100&p=" + page
				}
				fmt.Println(url)
				resp, err := http.Get(url)          // GETリクエストでアクセス
				if err != nil {
						return
				}
				defer resp.Body.Close()             // 関数が終了するとクローズ

				doc, err := goquery.NewDocumentFromResponse(resp) // レスポンスをNewDocumentFromResponseに渡してドキュメントを得る
				if err != nil {
						panic(err)
				}
				doc.Find("figure.pg-company-area-jobcassette-company").Each(func(_ int, s *goquery.Selection) {
						a := s.Find("a.huge")
						company := a.Text()
						onClick, _ := a.Attr("onclick")
						jumpUrl := rootUrl + ExtractUrlRegExp(onClick)
						address := s.Find("li.pg-company-area-jobcassette-company-address>p").Text()
						size := ""
						capital := ""
						s.Find("li.pg-company-area-jobcassette-company-capital>ul>li").Each(func(_ int, s *goquery.Selection) {
							str := s.Text()
							initials := string([]rune(str)[:2])
							if initials == "会社" {
									size = string([]rune(str)[5:])
							} else if initials == "資本" {
									capital = CapitalNumRegExp(string([]rune(str)[4:]))
							}
						})
						if len(size) == 0 {
								size = "0"
						}
						if len(capital) == 0 {
								capital = "0"
						}
						industry := ""
						s.Find("li.pg-company-area-jobcassette-company-industry>span").Each(func(_ int, s *goquery.Selection) {
								industry += s.Text() + "/"
						})
						file.Write(([]byte)(company+"，"+size+"，"+capital+"，"+industry+"，"+address+"，"+jumpUrl+"\n"))
				})

				count++                                                                // アクセスが成功したことをカウント
    }
    end := time.Now()                                                          // 処理にかかった時間を測定
    log.Printf("%d 回のリクエストに成功しました！\n", count)                        // 成功したリクエストの数を表示
    log.Printf("%f 秒処理に時間がかかりました！\n",(end.Sub(start)).Seconds())     //何秒かかったかを表示
}
