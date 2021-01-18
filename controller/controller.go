package controller

import (
	"context"
	"github.com/MaxFuhrich/containerNetworkExample/entities"
	"github.com/MaxFuhrich/containerNetworkExample/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

var collection *mongo.Collection
var ctx = context.TODO()

//Setting up connection to MongoDB
func init() {
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("timeentries").Collection("times")
}

func AddTime(c *gin.Context) {
	var t entities.RequestTime
	t = service.GetTime()
	_, err := collection.InsertOne(ctx, t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, t)
}

func History(c *gin.Context) {
	var entries []*entities.RequestTime
	filter := bson.D{}
	current, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	for current.Next(ctx) {
		var t entities.RequestTime
		err := current.Decode(&t)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		entries = append(entries, &t)
	}
	err = current.Close(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if len(entries) == 0 {
		c.JSON(http.StatusNoContent, mongo.ErrNoDocuments)
		return
	}
	c.JSON(http.StatusOK, entries)

}
