package gtw

type router struct {
	Handlers   map[string]HandlerFunc
	Middleware []Middleware
}
