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

// DeleteAchievement menghapus achievement dari MongoDB (hard delete)
func DeleteAchievement(id primitive.ObjectID) error {
	collection := config.GetMongoDB().Collection("achievements")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// SoftDeleteAchievement melakukan soft delete achievement di MongoDB
func SoftDeleteAchievement(id primitive.ObjectID) error {
	collection := config.GetMongoDB().Collection("achievements")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"deletedAt": now,
				"updatedAt": now,
			},
		},
	)

	return err
}

// GetAchievementsByMongoIDs mengambil multiple achievements berdasarkan array of MongoDB IDs
func GetAchievementsByMongoIDs(mongoIDs []string) ([]mongodb.Achievement, error) {
	collection := config.GetMongoDB().Collection("achievements")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert string IDs to ObjectIDs
	objectIDs := make([]primitive.ObjectID, 0, len(mongoIDs))
	for _, id := range mongoIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue // Skip invalid IDs
		}
		objectIDs = append(objectIDs, objectID)
	}

	// Query dengan $in operator
	filter := bson.M{
		"_id":       bson.M{"$in": objectIDs},
		"deletedAt": bson.M{"$exists": false}, // Exclude soft deleted
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []mongodb.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

// GetAchievementStatsByType mengambil statistik prestasi per tipe
func GetAchievementStatsByType(mongoIDs []string) (map[string]int, error) {
	collection := config.GetMongoDB().Collection("achievements")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert string IDs to ObjectIDs
	objectIDs := make([]primitive.ObjectID, 0, len(mongoIDs))
	for _, id := range mongoIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}

	// Aggregation pipeline
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id":       bson.M{"$in": objectIDs},
				"deletedAt": bson.M{"$exists": false},
			},
		},
		{
			"$group": bson.M{
				"_id":   "$achievementType",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	stats := make(map[string]int)
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		stats[result.ID] = result.Count
	}

	return stats, nil
}

// GetAchievementStatsByPeriod mengambil statistik prestasi per periode (bulan-tahun)
func GetAchievementStatsByPeriod(mongoIDs []string) (map[string]int, error) {
	collection := config.GetMongoDB().Collection("achievements")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert string IDs to ObjectIDs
	objectIDs := make([]primitive.ObjectID, 0, len(mongoIDs))
	for _, id := range mongoIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}

	// Aggregation pipeline - group by year-month
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id":       bson.M{"$in": objectIDs},
				"deletedAt": bson.M{"$exists": false},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"year":  bson.M{"$year": "$createdAt"},
					"month": bson.M{"$month": "$createdAt"},
				},
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"_id.year": -1, "_id.month": -1},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	stats := make(map[string]int)
	for cursor.Next(ctx) {
		var result struct {
			ID struct {
				Year  int `bson:"year"`
				Month int `bson:"month"`
			} `bson:"_id"`
			Count int `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		// Format: "2024-01"
		period := time.Date(result.ID.Year, time.Month(result.ID.Month), 1, 0, 0, 0, 0, time.UTC).Format("2006-01")
		stats[period] = result.Count
	}

	return stats, nil
}

// GetCompetitionLevelDistribution mengambil distribusi tingkat kompetisi
func GetCompetitionLevelDistribution(mongoIDs []string) (map[string]int, error) {
	collection := config.GetMongoDB().Collection("achievements")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert string IDs to ObjectIDs
	objectIDs := make([]primitive.ObjectID, 0, len(mongoIDs))
	for _, id := range mongoIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}

	// Aggregation pipeline - group by competition level
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id":             bson.M{"$in": objectIDs},
				"deletedAt":       bson.M{"$exists": false},
				"achievementType": "competition",
			},
		},
		{
			"$group": bson.M{
				"_id":   "$details.level",
				"count": bson.M{"$sum": 1},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	stats := make(map[string]int)
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		if result.ID == "" {
			result.ID = "unknown"
		}
		stats[result.ID] = result.Count
	}

	return stats, nil
}
