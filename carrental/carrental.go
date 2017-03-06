package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "encoding/json"
    "io"
    "io/ioutil"
    "math/rand"
    "encoding/csv"
    "os"
    "strconv"
    "bufio"
)

type Order struct {
    CarMaker string
    CarModel string
    NDays int
    NUnits int
    Price int
}


func main() {

router := mux.NewRouter().StrictSlash(true)
router.HandleFunc("/", Index)
router.HandleFunc("/rentals/new", handleNewOrder)
router.HandleFunc("/rentals/list", handleListOrders)

/*
router.HandleFunc("/endpoint/{param}", endpointFunc)
router.HandleFunc("/endpoint2/{param}", endpointFunc2JSONInput)
*/

log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Service OK")
}


func handleNewOrder(w http.ResponseWriter, r *http.Request) {
    var requestMessage Order
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &requestMessage); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    } else {
	requestMessage.Price = rand.Intn(100)*requestMessage.NDays*requestMessage.NUnits
        //res := Order{OrderData: requestMessage, Price: preu}
        writeOrderToFile(w,requestMessage)
        json.NewEncoder(w).Encode(requestMessage)
    }

}

func handleListOrders(w http.ResponseWriter, r *http.Request) {
    var lines []Order
    file, err := os.Open("rentals.csv")
    if err!=nil {
    json.NewEncoder(w).Encode(err)
    return
    }
    reader := csv.NewReader(bufio.NewReader(file))
    for {
        record, err := reader.Read()
        if err == io.EOF {
                break
        }
            //fmt.Fprintf(w, "The first value is %q", record[0])
        lines = append(lines, Order{CarMaker: record[0], CarModel: record[1],
        NDays: strconv.Atoi(record[2]), NUnits: strconv.Atoi(record[3]),
         Price: strconv.Atoi(record[4]) })
    }
        json.NewEncoder(w).Encode(lines)
}

func writeOrderToFile(w http.ResponseWriter, o Order) {
    file, err := os.OpenFile("rentals.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err!=nil {
        json.NewEncoder(w).Encode(err)
        return
    }
    writer := csv.NewWriter(file)
    var data1 = []string{o.CarMaker,o.CarModel,
    strconv.Itoa(o.NDays),strconv.Itoa(o.NUnits),strconv.Itoa(o.Price)}
    writer.Write(data1)
    writer.Flush()
    file.Close()
}

/*
func endpointFunc(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    param := vars["param"]
    res := Order{Field1: "Text1", Field2: param}
    json.NewEncoder(w).Encode(res)
}

func endpointFunc2JSONInput(w http.ResponseWriter, r *http.Request) {
    var requestMessage RequestMessage
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &requestMessage); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    } else {
        fmt.Fprintln(w, "Successfully received request with Field1 =", requestMessage.Field1)
        fmt.Println(r.FormValue("queryparam1"))
    }
}
*/
