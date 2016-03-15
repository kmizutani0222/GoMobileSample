package main

import (
    "fmt"
    "html"
    "net/http"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
    param1 := req.FormValue("param1")
    param2 := req.FormValue("param2")
    output := `{"header":{"code":200,"message":"OK"},"body":{"param1":"` + html.EscapeString(param1) + `", "param2":"` + html.EscapeString(param2) + `"}`
    fmt.Fprintf(w, "%s", output)
}

func main() {
    http.HandleFunc("/apisample", HelloServer)
    http.ListenAndServe(":8080", nil)
}
