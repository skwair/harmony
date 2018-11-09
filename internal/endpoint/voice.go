package endpoint

func GetVoiceRegions() *Endpoint {
	return &Endpoint{
		URL: "/voice/regions",
		Key: "/voice/regions",
	}
}
