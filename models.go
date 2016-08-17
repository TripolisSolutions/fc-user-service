package main

import "time"

type (
	IdentityCard struct {
		Number string `json:"number"`
	}

	Vehicle struct {
		VehicleType  string `json:"type"`
		LicensePlate string `json:"license_plate"`
	}

	User struct {
		ID           string       `json:"id"`
		Name         string       `json:"name"`
		EmailAddress string       `json:"email_address"`
		PhoneNumber  string       `json:"phone_number"`
		IdentityCard IdentityCard `json:"identity_card"`
		Vehicles     []Vehicle    `json:"vehicles"`
		CreatedAt    time.Time    `bson:"c_at" json:"c_at"`
		UpdatedAt    time.Time    `bson:"u_at" json:"u_at"`
	}
)
