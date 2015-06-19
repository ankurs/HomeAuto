package auto

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
	"log"
	"net/http"
	"strings"
)

const (
	RELAY_1 = 24
	RELAY_2 = 23
	RELAY_4 = 18
)

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("got request for index")
	fmt.Fprint(w, "Hello World")
}

func on(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sw, ok := vars["switch"]
	if ok {
		log.Println("got request for on " + sw)
		resp := Response{}
		status := http.StatusOK
		err := switchRelay(sw, embd.Low)
		if err == nil {
			resp.Status = "OK"
			resp.Msg = "On" + sw
		} else {
			resp.Status = "Error"
			resp.Msg = err.Error()
		}
		writeJSON(w, status, resp)
	} else {
		log.Println("got error request for on")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("error")
	}
}

func off(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sw, ok := vars["switch"]
	if ok {
		log.Println("got request for off " + sw)
		resp := Response{}
		status := http.StatusOK
		err := switchRelay(sw, embd.High)
		if err == nil {
			resp.Status = "OK"
			resp.Msg = "Off " + sw
		} else {
			resp.Status = "Error"
			resp.Msg = err.Error()
		}
		writeJSON(w, status, resp)
	} else {
		log.Println("got error request for off")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("error")
	}
}

func switchRelay(sw string, value int) error {
	embd.LEDToggle("LED0")
	switch strings.ToLower(strings.TrimSpace(sw)) {
	case "r1":
		return embd.DigitalWrite(RELAY_1, value)
	case "r2":
		return embd.DigitalWrite(RELAY_2, value)
	case "r4":
		return embd.DigitalWrite(RELAY_4, value)
	default:
		return errors.New("Relay " + sw + " not found")
	}
}

// writeJSON writes the value v to the http response stream as json with standard
// json encoding.
func writeJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

// Starts the main server
func Start() {
	defer recover()
	defer cleanup()
	intiEmbd()
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/on/{switch}", on)
	r.HandleFunc("/off/{switch}", off)
	err := http.ListenAndServe("0.0.0.0:4000", r)
	if err != nil {
		log.Fatal(err)
	}
}

func intiEmbd() {
	embd.InitGPIO()
	embd.SetDirection(RELAY_1, embd.Out)
	embd.DigitalWrite(RELAY_1, embd.High)
	embd.SetDirection(RELAY_2, embd.Out)
	embd.DigitalWrite(RELAY_2, embd.High)
	embd.SetDirection(RELAY_4, embd.Out)
	embd.DigitalWrite(RELAY_4, embd.High)
}

func cleanup() {
	embd.CloseGPIO()
}
