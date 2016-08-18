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

	user.UpdatedAt = time.Now()
	user.CreatedAt = time.Now()

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

	user.UpdatedAt = time.Now()

	if err := mongo.Execute("monotonic", UserCollection(tenantID),
		func(collection *mgo.Collection) error {
			return collection.UpdateId(user.ID, user)
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
				"id": user.ID,
			})
		}); err != nil {
		return err
	}

	return nil
}

// FindByID user
func (user *User) FindByID(tenantID string) error {
	if err := mongo.FindByID(user, UserCollection(tenantID), user.ID, nil); err != nil {
		return err
	}
	return nil
}

// FindUsers ...
func FindUsers(tenantID string) (users []User, err error) {
	//	if err := mongo.Execute("monotonic", UserCollection(tenantID),
	//		func(collection *mgo.Collection) error {
	//			return collection.Find(bson.M{
	//				"module": module,
	//			}).All(&categories)
	//		}); err != nil {
	//		return categories, err
	//	}

	return users, nil
}
