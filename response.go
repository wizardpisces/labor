package labor

import (
	"io/ioutil"
	"net/http/httptest"
	"syscall/js"
)

type ResponseRecorder struct {
	*httptest.ResponseRecorder
}

type Wrapper interface{ JSValue() js.Value }

func NewResponseRecorder() ResponseRecorder {
	return ResponseRecorder{httptest.NewRecorder()}
}

var _ Wrapper = ResponseRecorder{}

// var _ js.Wrapper = ResponseRecorder{}

func (rr ResponseRecorder) JSValue() js.Value {
	var res = rr.Result()

	var body js.Value = js.Undefined()
	if res.ContentLength != 0 {
		var b, err = ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		body = js.Global().Get("Uint8Array").New(len(b))
		js.CopyBytesToJS(body, b)
	}

	var init = make(map[string]interface{}, 2)

	if res.StatusCode != 0 {
		init["status"] = res.StatusCode
	}

	if len(res.Header) != 0 {
		var headers = make(map[string]interface{}, len(res.Header))
		for k := range res.Header {
			headers[k] = res.Header.Get(k)
		}
		init["headers"] = headers
	}

	return js.Global().Get("Response").New(body, init)
}
