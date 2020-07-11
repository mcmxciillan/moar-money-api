package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"labix.org/v2/mgo/bson"
)

type bill struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string        `json:"Name"`
	Amount    string        `json:"Amount"`
	Type      string        `json:"Type"`
	Frequency string        `json:"Frequency"`
}

type allBills []bill

// var bills = allBills{
// 	{
// 		ID:        "_",
// 		Name:      "_",
// 		Amount:    "_",
// 		Type:      "_",
// 		Frequency: "_",
// 	},
// }
var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
var client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
var collection = client.Database("moar").Collection("bills")

func createBill(w http.ResponseWriter, r *http.Request) {

	var newBill bill
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter correct Bill data")
	}

	json.Unmarshal(reqBody, &newBill)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, newBill)
	id := res.InsertedID
	fmt.Println("Succesfully inserted bill, ID: ", id)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newBill)
}

func getOneBill(w http.ResponseWriter, r *http.Request) {
	billID := bson.ObjectId(mux.Vars(r)["id"])
	fmt.Println("Getting here")
	expectedBill, err := findOneBill(billID)

	if err == nil {
		json.NewEncoder(w).Encode(expectedBill)
	}
}

func getAllBills(w http.ResponseWriter, r *http.Request) {

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Here's an array in which you can store the decoded documents
	var bills []*bill

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cur, err := collection.Find(ctx, bson.D{{}}, options.Find())
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem bill
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		bills = append(bills, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	cur.Close(context.TODO())
	json.NewEncoder(w).Encode(bills)
}

func updateBill(w http.ResponseWriter, r *http.Request) {

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	billID := bson.ObjectId(mux.Vars(r)["id"])
	_, err := findOneBill(billID)

	if err == nil {
		var updatedBill bill
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "Please enter correct bill data")
		}
		json.Unmarshal(reqBody, &updatedBill)

		updateResult, err := collection.UpdateOne(ctx, billID, updatedBill)

		if err != nil {
			log.Fatal(err)
			json.NewEncoder(w).Encode(err)
		}
		json.NewEncoder(w).Encode(updateResult)
	}
}

func deleteBill(w http.ResponseWriter, r *http.Request) {

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	billID := bson.ObjectId(mux.Vars(r)["id"])
	_, err := findOneBill(billID)

	deleteResult, err := collection.DeleteOne(ctx, billID)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request!"))
	}
	fmt.Println("Deleted bill: ", deleteResult)
}

// Finds a bill by it's ID
// Returns the expected bill if found, otherwise returns an error
func findOneBill(id bson.ObjectId) (bill, error) {

	// filter := {{"ID", id}}

	// The result we're expected to return
	var result bill

	// Get the database context
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Here")
	// Find the ID from the 'Bills' collection
	err = collection.FindOne(ctx, bson.D{{}}).Decode(&result)

	// If an error exists, log and return an error
	if err != nil {
		fmt.Println("OR Here")
		log.Fatal(err)
	}

	// Return the result
	return result, err
}
