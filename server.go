package gonoc

import (
    "net/http"
    "fmt"
    "log"
)

type Server struct {
}

func SetStaticRoute(path string, route string) {
    fsh := http.FileServer(http.Dir(path))
    http.Handle(route, http.StripPrefix(route, fsh))
}

func Run(path string, route string, listen string) {
    http.HandleFunc("/", Route)
    SetStaticRoute(path, route)
    errorHttp := http.ListenAndServe(listen, nil)
    if errorHttp != nil {
        log.Fatal("ListenAndServe :", errorHttp)
    }else{
        fmt.Println("Server OK")
    }
}