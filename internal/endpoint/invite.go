package endpoint

func GetInvite(code, query string) *Endpoint {
	return &Endpoint{
		URL: "/invites/" + code + "?" + query,
		Key: "/invites",
	}
}

func DeleteInvite(code string) *Endpoint {
	return &Endpoint{
		URL: "/invites/" + code,
		Key: "/invites",
	}
}
