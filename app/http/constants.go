package http

type ResponseStatus string
type HttpMethod string

const (
	GET  HttpMethod = "GET"
	POST HttpMethod = "POST"
)

const (
	OK        ResponseStatus = "200 OK"
	NOT_FOUND ResponseStatus = "404 Not Found"
	CREATED   ResponseStatus = "201 Created"
)

var supportedEncodings = []string{"gzip", "deflate"}
