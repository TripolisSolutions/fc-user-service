package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"goji.io/pat"
	"golang.org/x/net/context"

	"github.com/TripolisSolutions/go-helper/ids"
	"github.com/TripolisSolutions/go-helper/utilities"
)

type userPayload struct {
	CorrelationID string       `json:"correlation_id"`
	TenantID      string       `json:"tenant_id"`
	Name          string       `json:"name"`
	EmailAddress  string       `json:"email_address"`
	PhoneNumber   string       `json:"phone_number"`
	IdentityCard  IdentityCard `json:"identity_card"`
	Vehicles      []Vehicle    `json:"vehicles"`
}

func CreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tenantID := pat.Param(ctx, "tenant_id")

	log.WithFields(log.Fields{
		"tenantID": tenantID,
	}).Fatal("Requesting create user")

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.WithFields(log.Fields{
			"tenantID": tenantID,
			"error":    err,
		}).Fatal("Creating user: error while read request data")
		http.Error(w, "Error while reading request body", http.StatusBadRequest)
		return
	}

	var payl userPayload
	if err := json.Unmarshal(body, &payl); err != nil {
		log.WithFields(log.Fields{
			"tenantID": tenantID,
			"error":    err,
		}).Fatal("Creating user: error while unmarshalling request data")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	user := User{}
	user.ID, err = ids.GetID()
	if err != nil {
		log.WithFields(log.Fields{
			"tenantID": tenantID,
			"error":    err,
		}).Fatal("Creating user: error while get id from redis")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	user.Name = payl.Name
	user.EmailAddress = payl.EmailAddress
	user.PhoneNumber = payl.PhoneNumber
	user.IdentityCard = payl.IdentityCard
	user.Vehicles = payl.Vehicles

	if err := user.Insert(tenantID); err != nil {
		log.WithFields(log.Fields{
			"tenantID": tenantID,
			"user_id":  user.ID,
			"payload":  string(body),
			"error":    err,
		}).Fatal("Creating user: error while save user")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", utilities.ToJSON(user))
}

func GetUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	//	tenantID := pat.Param(ctx, "tenant_id")
	//	userID := pat.Param(ctx, "user_id")

	//	log.WithFields(log.Fields{
	//		"tenantID": tenantID,
	//		"userID":   userID,
	//	}).Fatal("Get user by tenant and id")
}

func UpdateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
}

func DeleteUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
}

func readRequestData(req *http.Request) (*User, error) {
	requestData := &User{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(requestData)
	if err != nil {
		//log
		return requestData, err
	}
	return requestData, nil
}
