package labor

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"syscall/js"
)

func Log(params ...interface{}) {
	js.Global().Get("console").Call("log", params[0])
}

func Request(r js.Value) *http.Request {
	jsBody := js.Global().Get("Uint8Array").New(Await(r.Call("arrayBuffer")))
	Log(jsBody)
	body := make([]byte, jsBody.Get("length").Int())
	js.CopyBytesToGo(body, jsBody)
	Log(2, body)
	req := httptest.NewRequest(
		r.Get("method").String(),
		r.Get("url").String(),
		bytes.NewBuffer(body),
	)
	Log(3)
	headersIt := r.Get("headers").Call("entries")
	for {
		e := headersIt.Call("next")
		if e.Get("done").Bool() {
			break
		}
		v := e.Get("value")
		req.Header.Set(v.Index(0).String(), v.Index(1).String())
	}
	Log(4, headersIt)
	return req
}
