package github

type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type GraphQLErrorResponse struct {
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}
