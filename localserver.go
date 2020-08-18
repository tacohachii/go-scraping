package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World")      // "Hello, World"と表示
}

func main() {
    http.HandleFunc("/", handler)       // http://localhost:8080/にアクセス -> handler関数を実行
    http.ListenAndServe(":8080", nil)   // サーバーを起動（ポート8080）
}
