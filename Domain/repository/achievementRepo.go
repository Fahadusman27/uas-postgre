package repository

import (
	"GOLANG/Domain/config"
	mongodb "GOLANG/Domain/model/mongoDB"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateAchievement menyimpan achievement ke MongoDB
func CreateAchievement(achievement *mongodb.Achievement) (*mongodb.Achievement, error) {
	collection := config.GetMongoDB().Collection("achievements")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set timestamps
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	result, err := collection.InsertOne(ctx, achievement)
	if err != nil {
		return nil, err
	}

	achievement.ID = result.InsertedID.(primitive.ObjectID)
	return achievement, nil
}

// GetAchievementByID mengambil achievement berdasarkan ID
func GetAchievementByID(id primitive.ObjectID) (*mongodb.Achievement, error) {
	collection := config.GetMongoDB().Collection("achievements")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var achievement mongodb.Achievement
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&achievement)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &achievement, nil
}

// UpdateAchievement update achievement di MongoDB
func UpdateAchievement(id primitive.ObjectID, achievement *mongodb.Achievement) error {
	collection := config.GetMongoDB().Collection("achievements")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	achievement.UpdatedAt = time.Now()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": achievement},
	)

	return err
}

// DeleteAchievement menghapus achievement dari MongoDB
func DeleteAchievement(id primitive.ObjectID) error {
	collection := config.GetMongoDB().Collection("achievements")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
