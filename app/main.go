package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/matoous/go-nanoid"
)

var redisPool *redis.Pool

func generateKey(str string) string {
	hash, err := gonanoid.Nanoid(7)
	if err != nil {
		panic(err)
	}
	return hash + "/" + str
}

func newKey(w http.ResponseWriter, r *http.Request) {
	var key string

	redisConn := redisPool.Get()
	defer redisConn.Close()

	vars := mux.Vars(r)
	keyName := vars["keyName"] // todo: check max len

	for {
		key = generateKey(keyName)
		hashExists, err := redis.Bool(redisConn.Do("EXISTS", key))
		if err != nil {
			log.Fatal(err)
		}
		if hashExists == false {
			break
		}
	}

	_, err := redisConn.Do("SET", key, "")
	if err != nil {
		log.Fatal(err)
	}

	rootDomain := "http://" + r.Host // is this safe?, also make http/https configurable
	newUrl := rootDomain + "/" + key

	fmt.Fprintf(w, "%s\n", newUrl)
}

func getKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]       // todo: check max len
	keyName := vars["keyName"] // todo: check max len

	// http.Error(w, http.StatusText(422), 422)

	key := hash + "/" + keyName

	redisConn := redisPool.Get()
	defer redisConn.Close()

	res, err := redis.String(redisConn.Do("GET", key))
	if err != nil {
		log.Println(err)
	}

	fmt.Fprintf(w, "%s", res)
}

func setKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]       // todo: check max len
	keyName := vars["keyName"] // todo: check max len
	value := vars["value"]     // todo: check max len

	// http.Error(w, http.StatusText(422), 422)

	key := hash + "/" + keyName

	redisConn := redisPool.Get()
	defer redisConn.Close()

	_, err := redis.String(redisConn.Do("SET", key, value))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Println("Starting server...")

	redisPool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "redis:6379")
		},
	}

	r := mux.NewRouter()

	r.HandleFunc("/new/{keyName}", newKey).Methods("POST")
	r.HandleFunc("/{hash}/{keyName}", getKey).Methods("GET")
	r.HandleFunc("/{hash}/{keyName}/{value}", setKey).Methods("POST")
	//r.HandleFunc("/{hash}/{keyName}", setKey).Methods("POST") // todo: add support for setting value in post body

	http.ListenAndServe(":80", r)
}
