// Code generated by goctl. DO NOT EDIT!
package model

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var prefixUserCacheKey = "cache:user:"

type userModel interface {
	Insert(ctx context.Context, data *User) error
	FindOne(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, data *User) error
	Delete(ctx context.Context, id string) error
}

type defaultUserModel struct {
	conn *monc.Model
}

func newDefaultUserModel(conn *monc.Model) *defaultUserModel {
	return &defaultUserModel{conn: conn}
}

func (m *defaultUserModel) Insert(ctx context.Context, data *User) error {
	if !data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}

	key := prefixUserCacheKey + data.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, data)
	return err
}

func (m *defaultUserModel) FindOne(ctx context.Context, id string) (*User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidObjectId
	}

	var data User
	key := prefixUserCacheKey + data.ID.Hex()
	err = m.conn.FindOne(ctx, key, &data, bson.M{"_id": oid})
	switch err {
	case nil:
		return &data, nil
	case monc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) Update(ctx context.Context, data *User) error {
	data.UpdateAt = time.Now()
	key := prefixUserCacheKey + data.ID.Hex()
	_, err := m.conn.ReplaceOne(ctx, key, bson.M{"_id": data.ID}, data)
	return err
}

func (m *defaultUserModel) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidObjectId
	}
	key := prefixUserCacheKey + id
	_, err = m.conn.DeleteOne(ctx, key, bson.M{"_id": oid})
	return err
}
