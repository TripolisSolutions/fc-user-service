package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/linkosmos/mapop"
	"goji.io/pat"
	"golang.org/x/net/context"

	"github.com/TripolisSolutions/go-helper/ids"
	logHelper "github.com/TripolisSolutions/go-helper/log"
	"github.com/TripolisSolutions/go-helper/utilities"
)

type userPayload struct {
	CorrelationID string       `json:"correlation_id"`
	TenantID      string       `json:"tenant_id"`
	Name          string       `json:"name" bson:"name"`
	EmailAddress  string       `json:"email_address" bson:"email_address"`
	PhoneNumber   string       `json:"phone_number" bson:"phone_number"`
	IdentityCard  IdentityCard `json:"identity_card" bson:"identity_card"`
	Vehicles      []Vehicle    `json:"vehicles" bson:"vehicles"`
}

func CreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tenantID := pat.Param(ctx, "tenant_id")

	logFields := log.Fields{
		"correlation_id": r.Header.Get(logHelper.HeaderCorrelationID),
		"tenantID":       tenantID,
		"host":           r.Host,
		"action":         "Creating user",
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.WithFields(mapop.Merge(logFields, log.Fields{
			"error": err,
		})).Fatal("Error while reading request data")
		http.Error(w, "Error while reading request data", http.StatusBadRequest)
		return
	}

	var payl userPayload
	if err := json.Unmarshal(body, &payl); err != nil {
		log.WithFields(mapop.Merge(logFields, log.Fields{
			"error": err,
		})).Fatal("Error while unmarshalling request data")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	user := User{}
	user.ID, err = ids.GetID()
	if err != nil {
		log.WithFields(mapop.Merge(logFields, log.Fields{
			"error": err,
		})).Fatal("Error while getting id from redis")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	user.Name = payl.Name
	user.EmailAddress = payl.EmailAddress
	user.PhoneNumber = payl.PhoneNumber
	user.IdentityCard = payl.IdentityCard
	user.Vehicles = payl.Vehicles

	if err := user.Insert(tenantID); err != nil {
		log.WithFields(mapop.Merge(logFields, log.Fields{
			"user_id": user.ID,
			"payload": string(body),
			"error":   err,
		})).Fatal("Error while saving user")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", utilities.ToJSON(user))
}

func UpdateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tenantID := pat.Param(ctx, "tenant_id")
	userID := pat.Param(ctx, "user_id")

	logFields := log.Fields{
		"correlation_id": r.Header.Get(logHelper.HeaderCorrelationID),
		"tenantID":       tenantID,
		"userID":         userID,
		"host":           r.Host,
		"action":         "Updating user",
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.WithFields(mapop.Merge(logFields, log.Fields{
			"error": err,
		})).Fatal("Error while reading request data")
		http.Error(w, "Error while reading request data", http.StatusBadRequest)
		return
	}

	var payl userPayload
	if err := json.Unmarshal(body, &payl); err != nil {
		log.WithFields(mapop.Merge(logFields, log.Fields{
			"error": err,
		})).Fatal("Error while unmarshalling request data")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	user := User{}
	user.ID = userID
	if err := user.FindByID(tenantID); err != nil {
		log.WithFields(mapop.Merge(logFields, log.Fields{
			"error": err,
		})).Fatal("Error while finding user by id")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	user.Name = payl.Name
	user.EmailAddress = payl.EmailAddress
	user.PhoneNumber = payl.PhoneNumber
	user.IdentityCard = payl.IdentityCard
	user.Vehicles = payl.Vehicles

	if err := user.Update(tenantID); err != nil {
		log.WithFields(mapop.Merge(logFields, log.Fields{
			"error": err,
		})).Fatal("Error while updating user")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", utilities.ToJSON(user))
}

func DeleteUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tenantID := pat.Param(ctx, "tenant_id")
	userID := pat.Param(ctx, "user_id")

	logFields := log.Fields{
		"correlation_id": r.Header.Get(logHelper.HeaderCorrelationID),
		"tenantID":       tenantID,
		"userID":         userID,
		"host":           r.Host,
		"action":         "Updating user",
	}

	user := User{}
	user.ID = userID
	if err := user.Delete(tenantID); err != nil {
		log.WithFields(mapop.Merge(logFields, log.Fields{
			"error": err,
		})).Fatal("Error while deleting user")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tenantID := pat.Param(ctx, "tenant_id")
	userID := pat.Param(ctx, "user_id")

	logFields := log.Fields{
		"correlation_id": r.Header.Get(logHelper.HeaderCorrelationID),
		"tenantID":       tenantID,
		"userID":         userID,
		"host":           r.Host,
		"action":         "Updating user",
	}

	user := User{}
	user.ID = userID
	if err := user.FindByID(tenantID); err != nil {
		log.WithFields(mapop.Merge(logFields, log.Fields{
			"error": err,
		})).Fatal("Error while finding user by id")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", utilities.ToJSON(user))
}

func GetUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	//	tenantID := pat.Param(ctx, "tenant_id")
	//	userID := pat.Param(ctx, "user_id")

	//	log.WithFields(log.Fields{
	//		"tenantID": tenantID,
	//		"userID":   userID,
	//	}).Fatal("Get user by tenant and id")
}
