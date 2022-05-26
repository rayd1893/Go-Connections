package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (r *Renderer) RenderJSON(w http.ResponseWriter, code int, data interface{}) {
	if data == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)

		if code >= 200 && code < 300 {
			fmt.Fprint(w, jsonOKResp)
			return
		}

		fmt.Fprintf(w, jsonErrTmpl, http.StatusText(code))
		return
	}

	if typ, ok := data.(error); ok {
		data = &singleError{Error: typ.Error()}
	}

	b := r.pool.Get().(*bytes.Buffer)
	b.Reset()
	defer r.pool.Put(b)

	if err := json.NewEncoder(b).Encode(data); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, jsonErrTmpl, http.StatusText(http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = b.WriteTo(w)
}

const jsonErrTmpl = `{"error":"%s"}`

const jsonOKResp = `{"ok":true}`

type singleError struct {
	Error string `json:"error,omitempty"`
}
