package main

import (
    "bytes"
    "crypto/md5"
    "encoding/binary"
    "encoding/json"
    "io"
    "net/http"
)

var list map[string]string

var charTable = [...]rune{
    'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k',
    'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
    'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6',
    '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
    'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S',
    'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

func index(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "It works!")
}

func shorten(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    url, ok := r.Form["url"]
    result := map[string]string{}

    if ok && url[0] != "" {
        sumData := md5.Sum([]byte(url[0]))
        partUint := binary.BigEndian.Uint32(sumData[0 : 4])
        shortBuffer := &bytes.Buffer{}
        for i := 0; i < 6; i++ {
            shortBuffer.WriteRune(charTable[partUint % 62])
            partUint = partUint >> 5
        }
        short := "http://" + r.Host + "/" + shortBuffer.String()

        list[short] = url[0]
        result["short"] = short
    } else {
        result["err"] = "param url empty"
    }

    bytes, _ := json.Marshal(result)
    io.WriteString(w, string(bytes))
}

func original(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    short, ok := r.Form["short"]
    result := map[string]string{}

    if ok && short[0] != "" {
        original := list[short[0]]
        result["original"] = original
    } else {
        result["err"] = "param short empty"
    }

    bytes, _ := json.Marshal(result)
    io.WriteString(w, string(bytes))
}

func main() {
    list = map[string]string{}
    http.HandleFunc("/", index)
    http.HandleFunc("/shorten", shorten)
    http.HandleFunc("/original", original)
    http.ListenAndServe(":80", nil)
}
