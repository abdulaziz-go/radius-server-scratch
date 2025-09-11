package handlers

import (
	"layeh.com/radius"
)

func AccessHandler(w radius.ResponseWriter, r *radius.Request) {
	var code radius.Code

	w.Write(r.Response(code))
}
