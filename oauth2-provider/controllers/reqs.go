package controllers

type GetAuthReq struct {
	ClientID     string `form:"client_id"`
	RedirectURI  string `form:"redirect_uri"`
	ResponseType string `form:"response_type"`
	Scope        string `form:"scope"`
	State        string `form:"state"`
}

type PostAuthReq struct {
	ClientID     string `form:"client_id"`
	RedirectURI  string `form:"redirect_uri"`
	ResponseType string `form:"response_type"`
	Scope        string `form:"scope"`
	State        string `form:"state"`
	Email        string `form:"email"`
	Password     string `form:"password"`
}
