package sofifa

import (
	"fifacm-scout/internal/models"
	"fmt"
	"github.com/gocolly/colly"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func GetPlayersAndTeams(databaseId uint) ([]models.Player, map[int]models.Team, error) {
	urlReplacer := strings.NewReplacer("https://", "", "http://", "")
	baseUrlWithoutProtocol := urlReplacer.Replace(baseURL)

	log.Println(fmt.Sprint("Going to retrieve players for database id: ", databaseId))

	players := make([]models.Player, 0)
	teams := map[int]models.Team{}
	offset := 0
	for offset >= 0 {
		c := colly.NewCollector(colly.AllowedDomains(baseUrlWithoutProtocol))
		pageAmount := 0

		player := models.Player{}

		c.OnResponse(func(r *colly.Response) {
			if r.StatusCode != 200 {
				log.Println(fmt.Sprint("Got error response:", r.StatusCode, ": ", r.Body))
				return
			}
		})

		c.OnHTML(".list tr", func(element *colly.HTMLElement) {
			avatarImages := element.ChildAttr(".col-avatar img", "data-srcset")
			player.PhotoURL = avatarImages[strings.LastIndex(avatarImages, "https"):strings.Index(avatarImages, " 3x")]

			player.DBUpdateID = databaseId
			name := element.ChildText("a[href^='/player/']")
			fullname := element.ChildAttr("a[href^='/player/']", "aria-label")
			t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
			player.Name = name
			player.NameNormalized, _, _ = transform.String(t, name)
			player.FullName = fullname
			player.FullNameNormalized, _, _ = transform.String(t, fullname)
			player.Positions = ""
			element.ForEach(".col-name .pos", func(i int, element *colly.HTMLElement) {
				player.Positions += element.Text + " "
			})
			player.Positions = strings.TrimSpace(player.Positions)

			replacer := strings.NewReplacer("€", "", "cm", "", "kg", "")

			player.Nationality = element.ChildAttr(".col-name .flag", "title")
			player.Age = getNumberFromString(element.ChildText("td.col-ae"))
			player.Rating = getNumberFromString(element.ChildText("td.col-oa"))
			player.Potential = getNumberFromString(element.ChildText("td.col-pt"))
			teamImages := element.ChildAttr(".col-name img.team", "data-srcset")
			if len(teamImages) > 0 {
				teamId := getNumberFromString(element.ChildAttr("td.col-name .ellipsis a", "href"))
				player.TeamID = teamId
				_, containsKey := teams[teamId]
				if !containsKey {
					teams[teamId] = models.Team{
						TeamID:     getNumberFromString(element.ChildAttr("td.col-name .ellipsis a", "href")),
						Name:       element.ChildText("a[href^='/team/']"),
						LogoUrl:    teamImages[strings.LastIndex(teamImages, "https"):strings.Index(teamImages, " 3x")],
						DBUpdateID: databaseId,
					}
				}
			}
			player.PlayerId = getNumberFromString(element.ChildText("td.col-pi"))
			player.Length = getNumberFromString(replacer.Replace(element.ChildText("td.col-hi")))
			player.Weight = getNumberFromString(replacer.Replace(element.ChildText("td.col-wi")))
			player.PreferredFoot = element.ChildText("td.col-pf")
			player.Value = transformMoney(replacer.Replace(element.ChildText("td.col-vl")))
			player.Wage = transformMoney(replacer.Replace(element.ChildText("td.col-wg")))
			player.ReleaseClause = transformMoney(replacer.Replace(element.ChildText("td.col-rc")))
			player.WeakFoot = getNumberFromString(element.ChildText("td.col-wk"))
			player.SkillMoves = getNumberFromString(element.ChildText("td.col-sk"))
			player.InternationalReputation = getNumberFromString(element.ChildText("td.col-ir"))
			player.BodyType = element.ChildText("td.col-bt")
			player.WorkRateAttacking = element.ChildText("td.col-aw")
			player.WorkRateDefending = element.ChildText("td.col-dw")

			onLoan := element.ChildText(".col-name .sub .bp3-tag")
			if strings.ToUpper(onLoan) == "ON LOAN" {
				player.LoanedUntil, _ = time.Parse("Jan 2, 2006", strings.ReplaceAll(onLoan, " ON LOAN", ""))
			} else {
				joined := element.ChildText("td.col-jt")
				if len(joined) > 0 {
					player.JoinedTeam, _ = time.Parse("Jan 2, 2006", joined)
				}

				contract := element.ChildText(".col-name .sub")
				if contract == "Free" {
					// No contract
				} else {
					player.ContractUntil, _ = time.Parse("2006", contract[strings.Index(contract, "~ "):])
				}
			}

			player.MentalAttributes = models.MentalAttributes{
				Aggression:    getNumberFromString(element.ChildText("td.col-ar")),
				Composure:     getNumberFromString(element.ChildText("td.col-cm")),
				Interceptions: getNumberFromString(element.ChildText("td.col-in")),
				Positioning:   getNumberFromString(element.ChildText("td.col-po")),
				Vision:        getNumberFromString(element.ChildText("td.col-vi")),
			}

			player.TechnicalAttributes = models.TechnicalAttributes{
				BallControl:        getNumberFromString(element.ChildText("td.col-bl")),
				Crossing:           getNumberFromString(element.ChildText("td.col-cr")),
				Curve:              getNumberFromString(element.ChildText("td.col-cu")),
				DefensiveAwareness: getNumberFromString(element.ChildText("td.col-ma")),
				Dribbling:          getNumberFromString(element.ChildText("td.col-dr")),
				Finishing:          getNumberFromString(element.ChildText("td.col-fi")),
				FKAccuracy:         getNumberFromString(element.ChildText("td.col-fr")),
				HeadingAccuracy:    getNumberFromString(element.ChildText("td.col-he")),
				LongPassing:        getNumberFromString(element.ChildText("td.col-lo")),
				LongShots:          getNumberFromString(element.ChildText("td.col-ln")),
				Penalties:          getNumberFromString(element.ChildText("td.col-ln")),
				ShortPassing:       getNumberFromString(element.ChildText("td.col-sh")),
				ShotPower:          getNumberFromString(element.ChildText("td.col-so")),
				SlidingTackle:      getNumberFromString(element.ChildText("td.col-sl")),
				StandingTackle:     getNumberFromString(element.ChildText("td.col-sa")),
				Volleys:            getNumberFromString(element.ChildText("td.col-vo")),
			}

			player.PhysicalAttributes = models.PhysicalAttributes{
				Acceleration: getNumberFromString(element.ChildText("td.col-ac")),
				SprintSpeed:  getNumberFromString(element.ChildText("td.col-sp")),
				Agility:      getNumberFromString(element.ChildText("td.col-ag")),
				Reactions:    getNumberFromString(element.ChildText("td.col-re")),
				Balance:      getNumberFromString(element.ChildText("td.col-ba")),
				Jumping:      getNumberFromString(element.ChildText("td.col-ju")),
				Stamina:      getNumberFromString(element.ChildText("td.col-st")),
				Strength:     getNumberFromString(element.ChildText("td.col-sr")),
			}

			player.GoalkeepingAttributes = models.GoalkeepingAttributes{
				GKDiving:      getNumberFromString(element.ChildText("td.col-gd")),
				GKHandling:    getNumberFromString(element.ChildText("td.col-gh")),
				GKKicking:     getNumberFromString(element.ChildText("td.col-gc")),
				GKPositioning: getNumberFromString(element.ChildText("td.col-gp")),
				GKReflexes:    getNumberFromString(element.ChildText("td.col-gr")),
			}

			players = append(players, player)
			pageAmount++
		})

		c.OnScraped(func(r *colly.Response) {
			if pageAmount == 60 {
				log.Println("amount was exactly 60, increasing offset for next batch")
				offset = offset + 60
			} else {
				log.Println("amount was less than 60, that means end was reached")
				offset = -1
			}
		})

		url := fmt.Sprint(baseURL, "/players?showCol[]=pi&showCol[]=ae&showCol[]=hi&showCol[]=wi&showCol[]=pf&showCol[]=oa"+
			"&showCol[]=pt&showCol[]=bo&showCol[]=bp&showCol[]=gu&showCol[]=le&showCol[]=vl&showCol[]=wg&showCol[]=rc"+
			"&showCol[]=ta&showCol[]=cr&showCol[]=fi&showCol[]=he&showCol[]=sh&showCol[]=ts&showCol[]=dr&showCol[]=cu"+
			"&showCol[]=fr&showCol[]=lo&showCol[]=bl&showCol[]=to&showCol[]=sp&showCol[]=ag&showCol[]=re&showCol[]=ba"+
			"&showCol[]=tp&showCol[]=so&showCol[]=ju&showCol[]=sr&showCol[]=ln&showCol[]=te&showCol[]=ar&showCol[]=in"+
			"&showCol[]=po&showCol[]=vi&showCol[]=cm&showCol[]=td&showCol[]=ma&showCol[]=sa&showCol[]=sl&showCol[]=tg"+
			"&showCol[]=gd&showCol[]=gc&showCol[]=gp&showCol[]=gr&showCol[]=bs&showCol[]=wk&showCol[]=sk&showCol[]=aw"+
			"&showCol[]=ir&showCol[]=bt&showCol[]=pac&showCol[]=sho&showCol[]=pas&showCol[]=dri&showCol[]=def"+
			"&showCol[]=vo&showCol[]=ac&showCol[]=st&showCol[]=gh&showCol[]=pe&showCol[]=dw&showCol[]=phy&showCol[]=tt"+
			"&showCol[]=jt&sort=desc&offset=", offset, "&r=", databaseId)
		log.Println(fmt.Sprint("Calling url with offset: ", offset))
		err := c.Visit(url)
		if err != nil {
			log.Println("error: " + err.Error())
			offset = -1
		}
	}
	return players, teams, nil
}

func ScrapePlayer(playerId int) (models.Player, error) {
	replacer := strings.NewReplacer("https://", "", "http://", "")
	baseUrlWithoutProtocol := replacer.Replace(baseURL)
	c := colly.NewCollector(colly.AllowedDomains(baseUrlWithoutProtocol))

	player := models.Player{}

	c.OnResponse(func(r *colly.Response) {
		log.Println(fmt.Sprint("Got response:", r.StatusCode, ", response size:", len(r.Body)))
		if r.StatusCode == 302 {
			return
		} else {
			player.PlayerId = playerId
		}
	})

	// Photo
	c.OnHTML(".player img:not(.flag)", func(element *colly.HTMLElement) {
		images := element.Attr("data-srcset")
		player.PhotoURL = images[strings.LastIndex(images, "https"):strings.Index(images, " 3x")]
	})

	c.OnHTML("body", func(element *colly.HTMLElement) {
		element.ForEach(".block-quarter", func(i int, childElement *colly.HTMLElement) {
			if i == 0 { // Rating
				rating, _ := strconv.Atoi(childElement.ChildText(".p"))
				player.Rating = rating
			} else if i == 4 { // Profile
				childElement.ForEach("ul > li", func(j int, liElement *colly.HTMLElement) {
					if j == 0 {
						player.PreferredFoot = strings.ReplaceAll(liElement.Text, "Preferred Foot", "")
					} else if j == 3 {
						rep, _ := strconv.Atoi(liElement.Text[0:1])
						player.InternationalReputation = rep
					} else if j == 4 {
						value := strings.ReplaceAll(liElement.Text, "Work Rate", "")
						player.WorkRateAttacking = value[:strings.Index(value, "/")]
						player.WorkRateDefending = value[strings.Index(value, "/ ")+2:]
					} else if j == 5 {
						player.BodyType = strings.ReplaceAll(liElement.Text, "Body Type", "")
					} else if j == 6 {
						value := strings.ReplaceAll(liElement.Text, "Real Face", "")
						if value == "Yes" {
							player.RealFace = true
						} else {
							player.RealFace = false
						}
					} else if j == 7 {
						if strings.Index(liElement.Text, "Release Clause") != -1 {
							player.ReleaseClause = getNumberFromString(strings.ReplaceAll(liElement.Text, "Release Clause€", ""))
						}
					}
				})
			} else if i == 6 { // Club
				images := childElement.ChildAttr("img:not(.flag)", "data-srcset")
				if len(images) > 0 {
					player.Team = models.Team{
						TeamID:  getNumberFromString(childElement.ChildAttr("a", "href")),
						Name:    childElement.ChildText("h5"),
						LogoUrl: images[strings.LastIndex(images, "https"):strings.Index(images, " 3x")]}
				}
			} else if i == 7 { // Physical
				childElement.ForEach("li", func(j int, liElement *colly.HTMLElement) {
					if j == 0 {
						player.PhysicalAttributes.Acceleration = getNumberFromString(liElement.Text)
					} else if j == 1 {
						player.PhysicalAttributes.Agility = getNumberFromString(liElement.Text)
					} else if j == 2 {
						player.PhysicalAttributes.Balance = getNumberFromString(liElement.Text)
					} else if j == 3 {
						player.PhysicalAttributes.Jumping = getNumberFromString(liElement.Text)
					} else if j == 4 {
						player.PhysicalAttributes.Reactions = getNumberFromString(liElement.Text)
					} else if j == 5 {
						player.PhysicalAttributes.SprintSpeed = getNumberFromString(liElement.Text)
					} else if j == 6 {
						player.PhysicalAttributes.Stamina = getNumberFromString(liElement.Text)
					} else if j == 7 {
						player.PhysicalAttributes.Strength = getNumberFromString(liElement.Text)
					}
				})
			} else if i == 8 { // Mental
				childElement.ForEach("li", func(j int, liElement *colly.HTMLElement) {
					if j == 0 {
						player.MentalAttributes.Aggression = getNumberFromString(liElement.Text)
					} else if j == 1 {
						player.MentalAttributes.Positioning = getNumberFromString(liElement.Text)
					} else if j == 2 {
						player.MentalAttributes.Composure = getNumberFromString(liElement.Text)
					} else if j == 3 {
						player.MentalAttributes.Interceptions = getNumberFromString(liElement.Text)
					} else if j == 4 {
						player.MentalAttributes.Vision = getNumberFromString(liElement.Text)
					}
				})
			} else if i == 9 { // Technical
				childElement.ForEach("li", func(j int, liElement *colly.HTMLElement) {
					if j == 0 {
						player.TechnicalAttributes.BallControl = getNumberFromString(liElement.Text)
					} else if j == 1 {
						player.TechnicalAttributes.Crossing = getNumberFromString(liElement.Text)
					} else if j == 2 {
						player.TechnicalAttributes.Curve = getNumberFromString(liElement.Text)
					} else if j == 3 {
						player.TechnicalAttributes.DefensiveAwareness = getNumberFromString(liElement.Text)
					} else if j == 4 {
						player.TechnicalAttributes.Dribbling = getNumberFromString(liElement.Text)
					} else if j == 5 {
						player.TechnicalAttributes.FKAccuracy = getNumberFromString(liElement.Text)
					} else if j == 6 {
						player.TechnicalAttributes.Finishing = getNumberFromString(liElement.Text)
					} else if j == 7 {
						player.TechnicalAttributes.HeadingAccuracy = getNumberFromString(liElement.Text)
					} else if j == 8 {
						player.TechnicalAttributes.LongPassing = getNumberFromString(liElement.Text)
					} else if j == 9 {
						player.TechnicalAttributes.LongShots = getNumberFromString(liElement.Text)
					} else if j == 10 {
						player.TechnicalAttributes.Penalties = getNumberFromString(liElement.Text)
					} else if j == 11 {
						player.TechnicalAttributes.ShortPassing = getNumberFromString(liElement.Text)
					} else if j == 12 {
						player.TechnicalAttributes.ShotPower = getNumberFromString(liElement.Text)
					} else if j == 13 {
						player.TechnicalAttributes.SlidingTackle = getNumberFromString(liElement.Text)
					} else if j == 14 {
						player.TechnicalAttributes.StandingTackle = getNumberFromString(liElement.Text)
					} else if j == 15 {
						player.TechnicalAttributes.Volleys = getNumberFromString(liElement.Text)
					}
				})
			} else if i == 10 { // Goalkeeping
				childElement.ForEach("li", func(j int, liElement *colly.HTMLElement) {
					if j == 0 {
						player.GoalkeepingAttributes.GKDiving = getNumberFromString(liElement.Text)
					} else if j == 1 {
						player.GoalkeepingAttributes.GKHandling = getNumberFromString(liElement.Text)
					} else if j == 2 {
						player.GoalkeepingAttributes.GKKicking = getNumberFromString(liElement.Text)
					} else if j == 3 {
						player.GoalkeepingAttributes.GKPositioning = getNumberFromString(liElement.Text)
					} else if j == 4 {
						player.GoalkeepingAttributes.GKReflexes = getNumberFromString(liElement.Text)
					}
				})
			}
		})
	})

	log.Println(fmt.Sprint("scraping for playerid:", playerId))
	err := c.Visit(fmt.Sprint(baseURL, "/player/", playerId, "?attr=career"))
	if err != nil {
		log.Println("error: " + err.Error())
		return player, err
	}

	if player.PlayerId == 0 {
		return player, fmt.Errorf("could not find player with id: %q", playerId)
	}

	return player, nil
}

func transformMoney(s string) int {
	result := strings.TrimSpace(s)
	lastChar := result[len(s)-1:]
	replacer := strings.NewReplacer(".", "")
	result = result[:len(result)-1]
	result = replacer.Replace(result)
	if lastChar == "M" {
		result += "00000"
		if !strings.Contains(s, ".") {
			result += "0"
		}
	} else if lastChar == "K" {
		result += "00"
		if !strings.Contains(s, ".") {
			result += "0"
		}
	}
	return getNumberFromString(result)
}

func getNumberFromString(s string) int {
	regex := regexp.MustCompile(`([0-9])+`)
	intValue, _ := strconv.Atoi(regex.FindString(s))
	return intValue
}
