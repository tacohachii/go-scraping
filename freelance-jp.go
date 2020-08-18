package main

import (
		"fmt"
    "log"
    "time"
		"os"
		"strconv"
		"strings"
		"net/http"
		"github.com/PuerkitoBio/goquery"
)

func main() {
		rootUrl := "https://www.freelance-jp.org"
		url := rootUrl + "/talents"           // アクセスするURL

    file, err := os.Create(`freelance-jp.csv`)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    count := 0                              // いくつアクセスが成功したかをカウント
    start := time.Now()                     // 処理にかかった時間を測定
    for maxRequest := 1; maxRequest <= 29; maxRequest ++ {     // リクエストを送る
				page := strconv.Itoa(maxRequest)
				url = rootUrl + "/talents?order=desc&page=" + page
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
				doc.Find(".p-talent__target").Each(func(_ int, s *goquery.Selection) {
						path, _ := s.Attr("href")
						href := rootUrl + path
						resp, err := http.Get(href)      // GETリクエストでアクセス
						if err != nil {
								return
						}
						defer resp.Body.Close()             // 関数が終了するとクローズ

						doc, err := goquery.NewDocumentFromResponse(resp) // レスポンスをNewDocumentFromResponseに渡してドキュメントを得る
						if err != nil {
								panic(err)
						}
						doc.Find(".p-profile__column").Each(func(_ int, s *goquery.Selection) {
								name := s.Find(".c-sidebar-profile.u-hide-desktop .c-sidebar-profile__name").Text()
								status := s.Find(".detail-info .c-tag-status .status").Text()
								compatible := ""
								s.Find(".detail-title span").Each(func(i int, s *goquery.Selection) {
										compatible += s.Text() + " / "
								})
								tel := "未登録"
								email := "未登録"
								s.Find(".detail-info ul").First().Find("li").Each(func(_ int, s *goquery.Selection) {
									src, _  := s.Find("img").Attr("src")
									if strings.Contains(src, "photo") {
										tel = s.Text()
										tel = strings.TrimSpace(tel)
									}
									if strings.Contains(src, "email") {
										email = s.Text()
										email = strings.TrimSpace(email)
									}
								})
								// CSVに出力
								file.Write(([]byte)(name+"，"+status+"，"+compatible+"，"+tel+"，"+email+"，"+href+"\n"))
						})
				})

				count++                                                                // アクセスが成功したことをカウント
    }
    end := time.Now()                                                          // 処理にかかった時間を測定
    log.Printf("%d 回のリクエストに成功しました！\n", count)                        // 成功したリクエストの数を表示
    log.Printf("%f 秒処理に時間がかかりました！\n",(end.Sub(start)).Seconds())     //何秒かかったかを表示
}
