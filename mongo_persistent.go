package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v2"

	mongo "github.com/TripolisSolutions/go-helper/mgojuice"
)

// UserCollection is ..
func UserCollection(tenantID string) string {
	return fmt.Sprintf("%s_users", tenantID)
}

// Insert user
func (user User) Insert(tenantID string) error {
	if err := validator.Validate(user); err != nil {
		return err
	}

	user.updatedAt = time.Now()
	user.createdAt = time.Now()

	if err := mongo.Execute("monotonic", UserCollection(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Insert(user)
		}); err != nil {
		return err
	}
	return nil
}

// Update user
func (user User) Update(tenantID string) error {
	if err := validator.Validate(user); err != nil {
		return err
	}

	user.updatedAt = time.Now()

	if err := mongo.Execute("monotonic", UserCollection(tenantID),
		func(collection *mgo.Collection) error {
			return collection.UpdateId(user.id, user)
		}); err != nil {
		return err
	}

	return nil
}

// Delete user
func (user User) Delete(tenantID string) error {

	if err := mongo.Execute("monotonic", UserCollection(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Remove(bson.M{
				"id": user.id,
			})
		}); err != nil {
		return err
	}

	return nil
}

// FindByID user
func (user *User) FindByID(tenant string, id string) error {
	if err := mongo.FindByID(user, UserCollection(tenant), id, nil); err != nil {
		return err
	}
	return nil
}
