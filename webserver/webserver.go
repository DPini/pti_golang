package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {

router := mux.NewRouter().StrictSlash(true)
router.HandleFunc("/", Index)
router.HandleFunc("/endpoint/{param}", endpointFunc)

log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Service OK")
}

func endpointFunc(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    param := vars["param"]
    fmt.Fprintln(w, "You are calling with param:", param)
}
