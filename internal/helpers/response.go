package helpers

func ApplyDefaultArg(messages []string, defaultArg string) string {
	if len(messages) > 0 {
		defaultArg = messages[0]
	}
	return defaultArg
}
