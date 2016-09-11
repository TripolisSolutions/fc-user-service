package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/linkosmos/mapop"
	"goji.io/pat"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"

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

	if err := deleteUserByID(tenantID, userID); err != nil {
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
	tenantID := pat.Param(ctx, "tenant_id")
	queries := r.URL.Query()

	limit := ParseIntWithFallback(queries.Get("limit"), 20)
	var offset int

	offset = ParseIntWithFallback(queries.Get("offset"), 0)

	if queries.Get("page") != "" {
		page := ParseIntWithFallback(queries.Get("page"), 0)
		offset = page * limit
	}

	if limit > 100 {
		limit = 100
	}

	var filterers = bson.M{}

	phoneNumber := strings.TrimSpace(queries.Get("phone_number"))
	if phoneNumber != "" {
		filterers["phone_number"] = phoneNumber
	}

	emailAddress := strings.TrimSpace(queries.Get("email_address"))
	if emailAddress != "" {
		filterers["email_address"] = emailAddress
	}

	identityCardNumber := strings.TrimSpace(queries.Get("identity_card_number"))
	if identityCardNumber != "" {
		filterers["identity_card.number"] = identityCardNumber
	}

	vehiclesLicensePlate := strings.TrimSpace(queries.Get("vehicles_license_plate"))
	if vehiclesLicensePlate != "" {
		filterers["vehicles.license_plate"] = vehiclesLicensePlate
	}

	log.WithFields(log.Fields{
		"filterers": string(utilities.ToJSON(filterers)),
		"limit":     limit,
		"offset":    offset,
	}).Info("get users properties")

	users, err := FindUsers(tenantID, filterers, limit, offset)
	if err != nil {
		log.WithFields(log.Fields{
			"filterers": filterers,
			"limit":     limit,
			"offset":    offset,
			"error":     err,
		}).Errorln("Failed to find users")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	total, err := CountUsers(tenantID, filterers)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Errorln("Failed to count users")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", utilities.ToJSON(struct {
		Docs  []User `json:"docs"`
		Total int    `json:"total"`
	}{
		Docs:  users,
		Total: total,
	}))
}
