package controllers

import (
	"fifacm-scout/internal/clients/sqlite"
	"fifacm-scout/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"strconv"
	"strings"
)

func GetPlayer(c *gin.Context) {
	playerId, _ := strconv.Atoi(c.Param("id"))

	player := models.Player{PlayerId: playerId}
	sqlite.DB.Preload(clause.Associations).Where(&player).First(&player)

	if player.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"message": "player not found"})
		return
	} else {
		c.JSON(200, player)
		return
	}

}

func GetPlayers(c *gin.Context) {
	// Query options
	limitParam := c.Query("limit")
	skipParam := c.Query("skip")
	sort := c.Query("sort")
	sortDirection := strings.ToLower(c.Query("sortDir"))

	limit := 200
	skip := 0
	if len(limitParam) > 0 {
		limit, _ = strconv.Atoi(limitParam)
	}
	if len(skipParam) > 0 {
		skip, _ = strconv.Atoi(skipParam)
	}

	if len(sort) < 1 {
		sort = "rating"
	}

	if sortDirection != "asc" && sortDirection != "desc" {
		sortDirection = "desc"
	}

	// Parameters
	dbId := c.Query("databaseId")

	databaseIdObject := models.DBUpdate{}
	sqlite.DB.Where(models.DBUpdate{UpdateID: dbId}).First(&databaseIdObject)

	db := sqlite.DB.
		Preload(clause.Associations).
		Debug().
		Where(models.Player{DBUpdateID: databaseIdObject.ID}).
		Order(sort + " " + sortDirection).
		Limit(limit).
		Offset(skip)

	if databaseIdObject.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"message": "databaseId not found"})
		return
	}

	addTermQuery(c, db)
	for _, attr := range []string{"age", "rating", "potential", "value", "wage", "length", "weight", "release_clause",
		"international_reputation", "weak_foot", "skill_moves"} {
		addBetweenQuery(c, db, attr)
	}
	addLikeQuery(c, db, "nationality")
	for _, physicalAttr := range []string{"acceleration", "agility", "balance", "jumping", "reactions", "sprint_speed", "stamina", "strength"} {
		addAssociationBetweenQuery(c, db, "physical_attributes", physicalAttr)
	}
	for _, gkAttr := range []string{"gk_diving", "gk_handling", "gk_kicking", "gk_positioning", "gk_reflexes"} {
		addAssociationBetweenQuery(c, db, "goalkeeping_attributes", gkAttr)
	}
	for _, mentalAttr := range []string{"agression", "positioning", "composure", "interceptions", "vision"} {
		addAssociationBetweenQuery(c, db, "mental_attributes", mentalAttr)
	}
	for _, technicalAttr := range []string{"ball_control", "crossing", "curve", "defensive_awareness", "dribbling", "fk_accuracy",
		"finishing", "heading_accuracy", "long_passing", "long_shots", "penalties", "short_passing", "shot_power", "sliding_tackle",
		"standing_tackle", "volleys"} {
		addAssociationBetweenQuery(c, db, "technical_attributes", technicalAttr)
	}
	addOrLikeQuery(c, db, "positions")

	var players []models.Player
	db.Find(&players)

	c.JSON(200, players)
}

func addTermQuery(c *gin.Context, db *gorm.DB) {
	term := c.Query("term")
	if len(term) > 0 {
		searchTerm := "%" + term + "%"
		db.Where(
			sqlite.DB.
				Where("name LIKE ?", searchTerm).
				Or("full_name LIKE ?", searchTerm).
				Or("name_normalized LIKE ?", searchTerm).
				Or("full_name_normalized LIKE ?", searchTerm),
		)
	}
}

func addBetweenQuery(c *gin.Context, db *gorm.DB, columnName string) {
	from, fromError := strconv.Atoi(c.Query(columnName + "From"))
	to, toError := strconv.Atoi(c.Query(columnName + "To"))

	if fromError == nil {
		db.Where(columnName+" >= ?", from)
	}
	if toError == nil {
		db.Where(columnName+" <= ?", to)
	}
}

func addLikeQuery(c *gin.Context, db *gorm.DB, columnName string) {
	queryValue := c.Query(columnName)
	if len(queryValue) > 0 {
		db.Where(columnName+" LIKE ?", "%"+queryValue+"%")
	}
}

func addAssociationBetweenQuery(c *gin.Context, db *gorm.DB, association string, columnName string) {
	from, fromError := strconv.Atoi(c.Query(columnName + "From"))
	to, toError := strconv.Atoi(c.Query(columnName + "To"))

	if fromError == nil && toError == nil {
		db.
			Joins("JOIN "+association+" a ON a.id == players."+association+"_id").
			Where(columnName+" >= ? AND "+columnName+" <= ?", from, to)
	} else if fromError == nil {
		db.
			Joins("JOIN physical_attributes pa ON pa.id == players.physical_attributes_id").
			Where(columnName+" >= ?", from)
	} else if toError == nil {
		db.
			Joins("JOIN physical_attributes pa ON pa.id == players.physical_attributes_id").
			Where(columnName+" <= ?", to)
	}
}

func addOrLikeQuery(c *gin.Context, db *gorm.DB, columnName string) {
	queryValue := c.Query(columnName)
	valueArray := strings.Split(queryValue, ",")

	if len(valueArray) > 0 {
		query := "("
		for idx, option := range valueArray {
			log.Println("adding: " + option)
			if idx != 0 {
				query += " OR "
			}
			query += columnName + " LIKE \"%" + option + "%\""
		}
		query += ")"
		db.Where(query)
	}
}
