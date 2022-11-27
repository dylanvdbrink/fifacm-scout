package main

import (
	"fifacm-scout/internal/clients/sofifa"
	"flag"
	"fmt"
)

func scrapePlayer() {
	playerId := flag.Int("id", 0, "The player id")
	flag.Parse()

	player, err := sofifa.ScrapePlayer(*playerId)
	if err != nil {
		panic(fmt.Sprint("Could not get player by id: '", *playerId, "': ", err.Error()))
	}

	fmt.Printf("%+v\n", player)
}
