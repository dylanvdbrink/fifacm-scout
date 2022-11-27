package main

import (
	"fifacm-scout/internal/clients/sofifa"
	"fifacm-scout/internal/clients/sqlite"
	"fifacm-scout/internal/models"
	"fmt"
	"log"
)

func main() {
	sqlite.Initialize()
	latestRemote := sofifa.GetLastDBUpdate()
	if isLatestDBUpdateSynced(latestRemote.UpdateID) {
		log.Println("latest local db updateid equals remote db updateid")
	} else {
		log.Println("local db not up to date, fetching players")
		getPlayers(latestRemote)
	}
}

func getPlayers(latestRemote models.DBUpdate) {
	sqlite.DB.Create(&latestRemote)

	// Get all the playerIds
	players, teams, err := sofifa.GetPlayersAndTeams(latestRemote.ID)
	if err != nil {
		log.Println("error: " + err.Error())
	}

	// Put players in database
	log.Println("creating players")
	for idx, player := range players {
		log.Println(fmt.Sprint("creating player with index: ", idx, ", playerid: ", player.PlayerId, " and name: ", player.Name))
		sqlite.DB.Create(&player)
	}

	// Put teams in database
	log.Println("creating teams")
	for idx, team := range teams {
		log.Println(fmt.Sprint("creating team with index: ", idx, ", teamid: ", team.TeamID, " and name: ", team.Name))
		sqlite.DB.Create(&team)
	}
}

func isLatestDBUpdateSynced(latestRemote string) bool {
	var dbUpdates []models.DBUpdate
	sqlite.DB.Order("update_id DESC").Find(&dbUpdates)
	if len(dbUpdates) == 0 {
		log.Println("no existing dbupdate yet, so not synced yet")
		return false
	}

	return dbUpdates[0].UpdateID == latestRemote
}
