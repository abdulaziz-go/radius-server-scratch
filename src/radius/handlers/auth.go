package handlers

import (
	"radius-server/src/database"
	cryptoUtil "radius-server/src/utils/crypto"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func AccessHandler(w radius.ResponseWriter, r *radius.Request) {
	username := rfc2865.UserName_GetString(r.Packet)
	password := rfc2865.UserPassword_GetString(r.Packet)

	user, err := database.GetUserByUsername(username)
	if err != nil || user == nil {
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}

	if !cryptoUtil.ComparePassword(user.PasswordHash, password) {
		w.Write(r.Response(radius.CodeAccessReject))
		return
	}

	w.Write(r.Response(radius.CodeAccessAccept))
}
