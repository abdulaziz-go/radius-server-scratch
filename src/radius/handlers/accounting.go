package handlers

import (
	"layeh.com/radius"
)

func AccountingHandler(w radius.ResponseWriter, r *radius.Request) {

	w.Write(r.Response(radius.CodeAccountingResponse))
}
