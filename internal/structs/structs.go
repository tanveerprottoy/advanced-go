package structs

func AnonymousStruct() any {
	return struct {
		message string
	}{
		"a message",
	}
}

