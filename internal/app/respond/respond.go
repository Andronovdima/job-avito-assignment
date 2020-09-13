package respond

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

func Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	if code == http.StatusInternalServerError {
		err = errors.New("We are sorry, some internal error on the server")
	}

	Respond(w, r, code, map[string]string{"message": errors.Cause(err).Error()})
}

func Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}
