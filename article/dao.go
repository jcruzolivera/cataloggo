package article

import (
	"context"

	"github.com/golang/glog"
	"github.com/nmarsollier/cataloggo/tools/db"
	"github.com/nmarsollier/cataloggo/tools/errs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrID = errs.NewValidation().Add("id", "Invalid")

// Define mongo Collection
var collection *mongo.Collection

func dbCollection() (*mongo.Collection, error) {
	if collection != nil {
		return collection, nil
	}

	database, err := db.Get()
	if err != nil {
		glog.Error(err)

		return nil, err
	}

	col := database.Collection("articles")

	collection = col
	return collection, nil
}

func findByCriteria(criteria string) ([]*Article, error) {
	var collection, err = dbCollection()
	if err != nil {
		glog.Error(err)

		return nil, err
	}

	filter := DBCriteriaFilter{
		Or: []map[string]DBCriteriaElement{
			{"description.name": {RegEx: criteria, Options: "i"}},
			{"description.description": {RegEx: criteria, Options: "i"}},
		},
	}

	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		glog.Error(err)

		return nil, err
	}
	defer cur.Close(context.Background())

	result := []*Article{}
	for cur.Next(context.Background()) {
		article := &Article{}
		if err := cur.Decode(article); err != nil {
			glog.Error(err)

			return nil, err
		}
		result = append(result, article)
	}

	return result, nil
}

type DBCriteriaFilter struct {
	Or []map[string]DBCriteriaElement `bson:"$or"`
}

type DBCriteriaElement struct {
	RegEx   string `bson:"$regex"`
	Options string `bson:"$options"`
}

func findById(articleId string) (*Article, error) {
	var collection, err = dbCollection()
	if err != nil {
		glog.Error(err)

		return nil, err
	}

	_id, err := primitive.ObjectIDFromHex(articleId)
	if err != nil {
		glog.Error(err)

		return nil, ErrID
	}

	article := &Article{}
	filter := DbIdFilter{ID: _id}

	if err = collection.FindOne(context.Background(), filter).Decode(article); err != nil {
		glog.Error(err)

		return nil, err
	}

	return article, nil
}

func insert(article *Article) (*Article, error) {
	if err := article.ValidateSchema(); err != nil {
		glog.Error(err)

		return nil, err
	}

	var collection, err = dbCollection()
	if err != nil {
		glog.Error(err)

		return nil, err
	}

	if _, err = collection.InsertOne(context.Background(), article); err != nil {
		glog.Error(err)

		return nil, err
	}

	return article, nil
}

// disable Deshabilita el articulo para que no se pueda usar mas
func Disable(articleId string) error {
	var collection, err = dbCollection()
	if err != nil {
		glog.Error(err)

		return err
	}

	_id, err := primitive.ObjectIDFromHex(articleId)
	if err != nil {
		glog.Error(err)

		return ErrID
	}

	_, err = collection.UpdateOne(context.Background(),
		DbIdFilter{ID: _id},
		DbEnableDocument{
			Set: DbEnableBody{
				Enabled: false,
			},
		},
	)

	return err
}

type DbEnableDocument struct {
	Set DbEnableBody `bson:"$set"`
}

type DbEnableBody struct {
	Enabled bool `bson:"enabled" json:"enabled"`
}

// Actualiza la descripción de un articulo.
func updateDescription(articleId string, description Description) error {
	var collection, err = dbCollection()
	if err != nil {
		glog.Error(err)

		return err
	}

	_id, err := primitive.ObjectIDFromHex(articleId)
	if err != nil {
		glog.Error(err)

		return ErrID
	}
	_, err = collection.UpdateOne(context.Background(),
		DbIdFilter{ID: _id},
		DbUpdateDescriptionDocument{
			Set: DbUpdateDescriptionBody{
				Description: description,
			},
		},
	)

	return err
}

type DbUpdateDescriptionDocument struct {
	Set DbUpdateDescriptionBody `bson:"$set"`
}

type DbUpdateDescriptionBody struct {
	Description Description `bson:"description"  json:"description" validate:"required"`
}

// Actualiza el precio de un articulo.
func updatePrice(articleId string, price float32) error {
	var collection, err = dbCollection()
	if err != nil {
		glog.Error(err)

		return err
	}

	_id, err := primitive.ObjectIDFromHex(articleId)
	if err != nil {
		glog.Error(err)

		return ErrID
	}
	_, err = collection.UpdateOne(context.Background(),
		DbIdFilter{ID: _id},
		DbUpdatePriceDocument{
			Set: DbUpdatePriceBody{
				Price: price,
			},
		},
	)

	return err
}

type DbUpdatePriceDocument struct {
	Set DbUpdatePriceBody `bson:"$set"`
}

type DbUpdatePriceBody struct {
	Price float32 `bson:"price"  json:"price"`
}

// Actualiza el stock de un articulo.
func updateStock(articleId string, stock int) error {
	var collection, err = dbCollection()
	if err != nil {
		glog.Error(err)

		return err
	}

	_id, err := primitive.ObjectIDFromHex(articleId)
	if err != nil {
		glog.Error(err)

		return ErrID
	}
	_, err = collection.UpdateOne(context.Background(),
		DbIdFilter{ID: _id},
		DbUpdateStockDocument{
			Set: DbUpdateStockBody{
				Stock: stock,
			},
		},
	)

	return err
}

type DbUpdateStockDocument struct {
	Set DbUpdateStockBody `bson:"$set"`
}

func DecreaseStock(articleId primitive.ObjectID, amount int) error {
	var collection, err = dbCollection()
	if err != nil {
		glog.Error(err)

		return err
	}

	_, err = collection.UpdateOne(context.Background(),
		DbIdFilter{ID: articleId},
		DbIncrementStockDocument{
			Set: DbIncrementStockBody{
				Inc: DbUpdateStockBody{
					Stock: -amount,
				},
			},
		},
	)

	return err
}

type DbIncrementStockBody struct {
	Inc DbUpdateStockBody `bson:"$inc"`
}

type DbIncrementStockDocument struct {
	Set DbIncrementStockBody `bson:"$set"`
}

type DbIdFilter struct {
	ID primitive.ObjectID `bson:"_id"`
}

type DbUpdateStockBody struct {
	Stock int `bson:"stock"  json:"stock"`
}
