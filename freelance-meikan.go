package main

import (
		"fmt"
    "log"
    "time"
		"os"
		"strconv"
		"net/http"
		"github.com/PuerkitoBio/goquery"
)

func main() {
		rootUrl := "https://freelance-meikan.com"
		url := rootUrl + "/freelance"           // アクセスするURL

    file, err := os.Create(`freelance-meikan.csv`)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    count := 0                              // いくつアクセスが成功したかをカウント
    start := time.Now()                     // 処理にかかった時間を測定
    for maxRequest := 1; maxRequest <= 32; maxRequest ++ {     // リクエストを送る
				if maxRequest != 1 {
					page := strconv.Itoa(maxRequest)
					url = rootUrl + "/freelance/page/" + page
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
				doc.Find(".freelancer-box").Each(func(_ int, s *goquery.Selection) {
						href, _ := s.Attr("href")
						resp, err := http.Get(href)      // GETリクエストでアクセス
						if err != nil {
								return
						}
						defer resp.Body.Close()             // 関数が終了するとクローズ

						doc, err := goquery.NewDocumentFromResponse(resp) // レスポンスをNewDocumentFromResponseに渡してドキュメントを得る
						if err != nil {
								panic(err)
						}
						doc.Find(".contents-body").Each(func(_ int, s *goquery.Selection) {
								name := s.Find("h1.name").Text()
								status := s.Find(".status-block dd").First().Text()
								compatible := ""
								occupation := ""
								s.Find(".fv .bottom .left dd").Each(func(i int, s *goquery.Selection) {
									if i == 0 {
											compatible = s.Find("span").Text()
									}
									if i == 1 {
										s.Find("span").Each(func(i int, s *goquery.Selection) {
												occupation += s.Text() + " / "
										})
									}
								})
								tel := ""
								email := ""
								s.Find("#contact .top a").Each(func(i int, s *goquery.Selection) {
									if i == 0 {
											tel = s.Find("dd span").Text()
									}
									if i == 1 {
											email = s.Find("dd span").Text()
									}
								})
								// CSVに出力
								file.Write(([]byte)(name+"，"+status+"，"+occupation+"，"+compatible+"，"+tel+"，"+email+"，"+href+"\n"))
						})
				})

				count++                                                                // アクセスが成功したことをカウント
    }
    end := time.Now()                                                          // 処理にかかった時間を測定
    log.Printf("%d 回のリクエストに成功しました！\n", count)                        // 成功したリクエストの数を表示
    log.Printf("%f 秒処理に時間がかかりました！\n",(end.Sub(start)).Seconds())     //何秒かかったかを表示
}
