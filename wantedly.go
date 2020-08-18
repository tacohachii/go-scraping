package main

import (
		"fmt"
    "log"
    "time"
		"os"
		"bufio"
		"strings"
		"net/http"
		"github.com/PuerkitoBio/goquery"
)

func getCompanyInfo(json string) (name string, hpUrl string, address string) {
	companyJson := strings.Split(strings.Split(json, "\"company\":{")[1], "},\"company_followings\"")[0]
	arr := strings.Split(companyJson, ",")
	name = replaceKey(arr[1], "name:")
	hpUrl = replaceKey(arr[4], "url:")
	address = replaceKey(arr[11], "short_location:") +
						replaceKey(arr[12], "address_prefix:") +
						replaceKey(arr[13], "address_suffix:")
	return name, hpUrl, address
}

func replaceKey(str string, key string) string {
	return strings.Replace(replaceDQ(str), key, "", 1)
}

func replaceDQ(str string) string {
	return strings.Replace(str, "\"", "", -1)
}

func main() {
		inFile, err := os.Open("google.csv")
		if err != nil {
			log.Fatalf("Error when opening file: %s", err)
		}
		fileScanner := bufio.NewScanner(inFile)

    outFile, err := os.Create(`wantedly.csv`)
    if err != nil {
        log.Fatal(err)
    }
    defer outFile.Close()

    count := 0                          // いくつアクセスが成功したかをアカウントするよ！
		start := time.Now()                 // 処理にかかった時間を測定
		
		for fileScanner.Scan() {            // 1行ずつファイルを読み込む
				url := fileScanner.Text()
				fmt.Println(url)
				resp, err := http.Get(url)      // GETリクエストでアクセス
				if err != nil {
						return
				}
				defer resp.Body.Close()         // 関数が終了するとクローズ

				doc, err := goquery.NewDocumentFromResponse(resp)  // レスポンスをNewDocumentFromResponseに渡してドキュメントを得る
				if err != nil {
						panic(err)
				}
				doc.Find("head > script").Each(func(_ int, s *goquery.Selection) {
						str := s.Text()
						if string([]rune(str)[:12]) == "// {\"router\"" {
								json := string([]rune(str)[3:])
								var company, hpUrl, address string
								company, hpUrl, address = getCompanyInfo(json)
								outFile.Write(([]byte)(company+"，"+hpUrl+"，"+address+"，"+url+"\n"))
						}
				})

				count++                         // アクセスが成功したことをカウント
		}
		if err := fileScanner.Err(); err != nil {
			log.Fatalf("Error while reading file: %s", err)
		}
		file.Close()

    end := time.Now()                                                         // 処理にかかった時間を測定
    log.Printf("%d 回のリクエストに成功しました！\n", count)                       // 成功したリクエストの数を表示
    log.Printf("%f 秒処理に時間がかかりました！\n",(end.Sub(start)).Seconds())     //何秒かかったかを表示
}
