package main

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"goji.io/pat"
	"golang.org/x/net/context"
)

type (
	IdentityCard struct {
		number string `json:"number"`
	}

	Vehicle struct {
		vehicleType  string `json:"type"`
		licensePlate string `json:"license_plate"`
	}

	User struct {
		id           string       `json:"id"`
		name         string       `json:"name"`
		emailAddress string       `json:"email_address"`
		phoneNumber  string       `json:"phone_number"`
		identityCard IdentityCard `json:"identity_card"`
		vehicles     []Vehicle    `json:"vehicles"`
		createdAt    time.Time    `bson:"c_at" json:"c_at"`
		updatedAt    time.Time    `bson:"u_at" json:"u_at"`
	}
)

func CreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {

}

func GetUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tenantUUID := pat.Param(ctx, "tenant_uuid")
	userID := pat.Param(ctx, "user_id")

	log.WithFields(log.Fields{
		"tenantUUID": tenantUUID,
		"userID":     userID,
	}).Fatal("Get user by tenant and id")
}

func readUserReqData(req *http.Request) (*User, error) {
	requestData := &User{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(requestData)
	if err != nil {
		//log
		return requestData, err
	}
	return requestData, nil
}
