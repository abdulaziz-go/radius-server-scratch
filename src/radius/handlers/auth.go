package handlers

import (
	"fmt"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func AccessHandler(w radius.ResponseWriter, r *radius.Request) {
	username := rfc2865.UserName_GetString(r.Packet)
	password := rfc2865.UserPassword_GetString(r.Packet)

	fmt.Printf("request %v %v", username, password)

	w.Write(r.Response(radius.CodeAccessAccept))
}
