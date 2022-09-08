package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"syscall/js"

	labor "github.com/yisar/labor"
)

func main() {

	js.Global().Get("labor").Get("http").Set("get", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		path := args[0].String()
		httpGet(path)
		return nil
	}))

	js.Global().Get("labor").Get("http").Set("serve", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		labor.Serve(nil)
		return nil
	}))

	js.Global().Get("labor").Get("http").Set("download", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		body := labor.Download(args[0].String())
		return body
	}))

	select {}
}
func Log(params ...interface{}) {
	js.Global().Get("console").Call("log", params[0])
}

func httpGet(path string) interface{} {
	http.HandleFunc(path, func(res http.ResponseWriter, req *http.Request) {
		params := make(map[string]string)
		if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
			fmt.Sprintf("params Hello %s!", params["name"])
			panic(err)
		}
		Log(4.1)
		res.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(map[string]string{
			"message": fmt.Sprintf("Hello %s!", params["name"]),
		}); err != nil {
			Log(4.2)
			panic(err)
		}
	})
	return nil
}
