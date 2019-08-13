package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/joho/godotenv"
)



type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")	

	resp, err := http.Post(os.Getenv("JAVA_MICROSERVICE_URL")+"/person", "application/json", request.Body)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	dd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var insertedID interface{}
	err = json.Unmarshal([]byte(dd), &insertedID)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("--------\n", insertedID)
    }
 
	result, err := json.Marshal(insertedID)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("--------\n", string(result))
    }

	json.NewEncoder(response).Encode(insertedID)

}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []Person
	resp, err := http.Get(os.Getenv("JAVA_MICROSERVICE_URL")+"/people")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body,&people)
	
	log.Println(string(body))
	json.NewEncoder(response).Encode(people)
	
}
func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	params := mux.Vars(request)
	log.Println(params["id"])
	resp, err := http.Get(os.Getenv("JAVA_MICROSERVICE_URL")+"/person/"+params["id"])
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body,&person)
	log.Println(string(body))
	json.NewEncoder(response).Encode(person)

}

func DeletePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	log.Println(params["id"])
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", os.Getenv("JAVA_MICROSERVICE_URL")+"/person/"+params["id"], nil)
    if err != nil {
        log.Fatalln(err)
        return
	}

	resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	var DeletedCount interface{}
	err = json.Unmarshal([]byte(body), &DeletedCount)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("--------\n", DeletedCount)
    }
 
	result, err := json.Marshal(DeletedCount)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("--------\n", string(result))
    }

	json.NewEncoder(response).Encode(DeletedCount)
}

func UpdatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", os.Getenv("JAVA_MICROSERVICE_URL")+"/person/"+params["id"], request.Body)
    if err != nil {
        log.Fatalln(err)
        return
	}

	resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

	var DeletedCount interface{}
	err = json.Unmarshal([]byte(body), &DeletedCount)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("--------\n", DeletedCount)
    }
 
	result, err := json.Marshal(DeletedCount)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("--------\n", string(result))
    }

	json.NewEncoder(response).Encode(DeletedCount)
}


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Starting the application...")
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	router := mux.NewRouter()
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", DeletePersonEndpoint).Methods("DELETE")
	router.HandleFunc("/person/{id}", UpdatePersonEndpoint).Methods("PUT")
	
	port := os.Getenv("PORT")
	fmt.Println("Listen at host"+port)

  log.Fatal(http.ListenAndServe(port, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(router)))
}