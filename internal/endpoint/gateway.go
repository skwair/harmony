package endpoint

func Gateway() *Endpoint {
	return &Endpoint{
		URL: "/gateway",
		Key: "/gateway",
	}
}

func GatewayBot() *Endpoint {
	return &Endpoint{
		URL: "/gateway/bot",
		Key: "/gateway/bot",
	}
}
