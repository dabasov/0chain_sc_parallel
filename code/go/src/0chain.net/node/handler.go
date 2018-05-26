package node

import "net/http"

func SetupHandlers() {
	http.HandleFunc("/_nh/whoami", WhoAmIHandler)
}

//WhoAmIHandler - who am i?
func WhoAmIHandler(w http.ResponseWriter, r *http.Request) {
	if Self == nil {
		return
	}
	Self.Print(w)
}