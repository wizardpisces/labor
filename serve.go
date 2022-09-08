package labor

import (
	"fmt"
	"net/http"
	"strings"
	"syscall/js"
)

func Serve(handler http.Handler) func() {
	var h = handler
	if h == nil {
		h = http.DefaultServeMux
	}

	var prefix = js.Global().Get("labor").Get("path").String()
	for strings.HasSuffix(prefix, "/") {
		prefix = strings.TrimSuffix(prefix, "/")
	}

	if prefix != "" {
		var mux = http.NewServeMux()
		mux.Handle(prefix+"/", http.StripPrefix(prefix, h))
		h = mux
	}

	var cb = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		var resPromise, resolve, reject = NewPromise()

		go func() {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						Log(8, err)
						reject(fmt.Sprintf("labor: panic: %+v\n", err))
					} else {
						reject(fmt.Sprintf("labor: panic: %v\n", r))
					}
				}
			}()

			var res = NewResponseRecorder()
			Log(args[0])
			Log(1, args, res)
			h.ServeHTTP(res, Request(args[0]))

			Log(5)
			resolve(res.JSValue())
			Log(7)

		}()

		return resPromise
	})

	js.Global().Get("labor").Call("setHandler", cb)

	return cb.Release
}
