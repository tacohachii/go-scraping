// go run test.go
// https://qiita.com/Azunyan1111/items/a1b6c58dc868814efb51
// https://qiita.com/akkikki_romeo/items/508331ef514da918aa8a
// https://qiita.com/Yaruki00/items/b50e346551690b158a79
// https://qiita.com/kou_pg_0131/items/dab4bcbb1df1271a17b6
// https://qiita.com/yosuke_furukawa/items/5fd41f5bcf53d0a69ca6
// https://qiita.com/derodero24/items/949ac666b18d567e9b61

package main

import (
		"strconv"
		"strings"
		"net/http"
		"sync"
    "log"
    "time"
		"os"
		"github.com/PuerkitoBio/goquery"
		// "io/ioutil"
    "fmt"
)

func ExtractUrlRegExp(str string) string {
	arr := strings.Split(str, "'")
	return arr[1]
}

func main() {
		// url := "http://localhost:8080"          // アクセスするURLだよ！
		// url := "https://www.bizreach.jp/company/"          // アクセスするURLだよ！
		rootUrl := "https://www.bizreach.jp"     
		url := rootUrl + "/company"       // アクセスするURLだよ！


    file, err := os.Create(`output.txt`)
    if err != nil {
        log.Fatal(err)  //ファイルが開けなかったときエラー出力
    }
    defer file.Close()

    maxConnection := make(chan bool,30)    // 同時に並列する数を指定できるよ！（第二引数）
    wg := &sync.WaitGroup{}                 // 並列処理が終わるまでSleepしてくれる便利なやつだよ！

    count := 0                              // いくつアクセスが成功したかをアカウントするよ！
    start := time.Now()                     // 処理にかかった時間を測定するよ！
    for maxRequest := 0; maxRequest < 28; maxRequest ++{     // 10000回リクエストを送るよ！
        wg.Add(1)                       // wg.add(1)とすると並列処理が一つ動いていることを便利な奴に教えるよ！
        maxConnection <- true               // ここは並列する数を抑制する奴だよ！詳しくはググって！
				go func() {                         // go func(){/*処理*/}とやると並列処理を開始してくれるよ！
            defer wg.Done()                 // wg.Done()を呼ぶと並列処理が一つ終わったことを便利な奴に教えるよ！

						if count != 0 {
							page := strconv.Itoa(count + 1)
							url = rootUrl + "/company/c/?pageSize=100&p=" + page
						}
						fmt.Println(url)
            resp, err := http.Get(url)      // GETリクエストでアクセスするよ！
            if err != nil {                 // err ってのはエラーの時にエラーの内容が入ってくるよ！
                return                      // 回線が狭かったりするとここでエラーが帰ってくるよ！
            }
						defer resp.Body.Close()         // 関数が終了するとなんかクローズするよ！（おまじない的な）

						// レスポンスをNewDocumentFromResponseに渡してドキュメントを得る
						doc, err := goquery.NewDocumentFromResponse(resp)
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
										size += string([]rune(str)[5:])
									} else if initials == "資本" {
										str = string([]rune(str)[4:])
										arr := strings.Split(strings.Split(str, "百万円")[0], ",")
										capital = strings.Join(arr, "")
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
								file.Write(([]byte)(company+"，"+jumpUrl+"，"+address+"，"+size+"，"+capital+"，"+industry+"\n"))
						})

						// HTML取得
						// byteArray, err := ioutil.ReadAll(resp.Body) // 帰ってきたレスポンスの中身を取り出すよ！
						// if err != nil {
						// 		panic(err)
						// }
						// file.Write(([]byte)(string(byteArray)+"\n"))

						// fmt.Println(string(byteArray))  // 取り出したリクエスト結果をバイナリ配列からstring型に変換して出力するよ！

            count++                         // アクセスが成功したことをカウントするよ！
            <-maxConnection                 // ここは並列する数を抑制する奴だよ！詳しくはググって！
        }()
    }
    wg.Wait()                               // ここは便利な奴が並列処理が終わるのを待つよ！
    end := time.Now()                       // 処理にかかった時間を測定するよ！
    log.Printf("%d 回のリクエストに成功しました！\n", count) // 成功したリクエストの数を表示してくれるよ！
    log.Printf("%f 秒処理に時間がかかりました！\n",(end.Sub(start)).Seconds())            //何秒かかったかを表示するよ！
}
