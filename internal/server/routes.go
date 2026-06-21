package server

// import "net/http"
//
// type Router struct {
// 	Base   string
// 	Routes []Route
// }
//
// type Route struct {
// 	Method  string
// 	Path    string
// 	Handler http.HandlerFunc
// }
//
// func NewRouter(base string) *Router {
// 	return &Router{
// 		Base:   base,
// 		Routes: nil,
// 	}
// }
//
// func NewRoute(method, path string, handler http.HandlerFunc) Route {
// 	return Route{
// 		Method:  method,
// 		Path:    path,
// 		Handler: handler,
// 	}
// }
//
// func (r *Router) RegisterNewRoute(method, path string, handler http.HandlerFunc) {
// 	route := NewRoute(method, path, handler)
// 	r.Routes = append(r.Routes, route)
// }
