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
func (user *User) Insert(tenantID string) error {
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
func (user *User) Update(tenantID string) error {
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
func deleteUserByID(tenantID, userID string) error {
	if err := mongo.Execute("monotonic", UserCollection(tenantID),
		func(collection *mgo.Collection) error {
			return collection.Remove(bson.M{
				"id": userID,
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

// FindUsers find users based on params
func FindUsers(tenantID string, filteres bson.M, limit, offset int) ([]User, error) {
	var users []User
	if err := mongo.Execute("monotonic", UserCollection(tenantID), func(collection *mgo.Collection) error {
		return collection.Find(filteres).Limit(limit).Skip(offset).Sort("-c_at").All(&users)
	}); err != nil {
		return users, err
	}

	return users, nil
}

// CountUsers count user based on params
func CountUsers(tenantID string, filteres bson.M) (int, error) {
	var result int
	if err := mongo.Execute("monotonic", UserCollection(tenantID), func(collection *mgo.Collection) error {
		count, err := collection.Find(filteres).Count()
		result = count
		return err
	}); err != nil {
		return result, err
	}

	return result, nil
}
