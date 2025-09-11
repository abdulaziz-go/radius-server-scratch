package handlers

import "layeh.com/radius"

func CoaHandler(w radius.ResponseWriter, r *radius.Request) {
	code := r.Packet.Code

	if code == radius.CodeDisconnectRequest {
		w.Write(r.Response(radius.CodeDisconnectACK))
		return
	}

	if code == radius.CodeCoARequest {
		w.Write(r.Response(radius.CodeCoAACK))
		return
	}

	w.Write(r.Response(radius.CodeCoANAK))
}
