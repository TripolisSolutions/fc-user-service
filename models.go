package main

import "time"

type (
	IdentityCard struct {
		Number string `json:"number" bson:"number"`
	}

	Vehicle struct {
		VehicleType  string `json:"type" bson:"type"`
		LicensePlate string `json:"license_plate" bson:"license_plate"`
	}

	User struct {
		ID           string       `json:"id" bson:"id"`
		Name         string       `json:"name" bson:"name"`
		EmailAddress string       `json:"email_address" bson:"email_address"`
		PhoneNumber  string       `json:"phone_number" bson:"phone_number"`
		IdentityCard IdentityCard `json:"identity_card" bson:"identity_card"`
		Vehicles     []Vehicle    `json:"vehicles" bson:"vehicles"`
		CreatedAt    time.Time    `json:"c_at" bson:"c_at"`
		UpdatedAt    time.Time    `json:"u_at" bson:"u_at"`
	}
)
