package models

import (
	"gorm.io/gorm"
	"time"
)

type Player struct {
	gorm.Model `json:"-"`

	PlayerId           int
	Name               string
	FullName           string
	PhotoURL           string
	Age                int
	Rating             int
	Potential          int
	Value              int
	Wage               int
	Length             int
	Weight             int
	Nationality        string
	FullNameNormalized string
	NameNormalized     string

	TeamID        int  `json:"-"`
	Team          Team `gorm:"references:TeamID"`
	JoinedTeam    time.Time
	ContractUntil time.Time
	LoanedUntil   time.Time

	Positions               string
	PreferredFoot           string
	WeakFoot                int
	SkillMoves              int
	InternationalReputation int
	WorkRateAttacking       string
	WorkRateDefending       string
	BodyType                string
	RealFace                bool
	ReleaseClause           int

	PhysicalAttributesID    int `json:"-"`
	PhysicalAttributes      PhysicalAttributes
	MentalAttributesID      int `json:"-"`
	MentalAttributes        MentalAttributes
	TechnicalAttributesID   int `json:"-"`
	TechnicalAttributes     TechnicalAttributes
	GoalkeepingAttributesID int `json:"-"`
	GoalkeepingAttributes   GoalkeepingAttributes

	DBUpdateID uint `json:"-"`
	DBUpdate   DBUpdate
}

type PhysicalAttributes struct {
	ID int `json:"-"`

	Acceleration int
	Agility      int
	Balance      int
	Jumping      int
	Reactions    int
	SprintSpeed  int
	Stamina      int
	Strength     int
}

type MentalAttributes struct {
	ID int `json:"-"`

	Aggression    int
	Positioning   int
	Composure     int
	Interceptions int
	Vision        int
}

type TechnicalAttributes struct {
	ID int `json:"-"`

	BallControl        int
	Crossing           int
	Curve              int
	DefensiveAwareness int
	Dribbling          int
	FKAccuracy         int
	Finishing          int
	HeadingAccuracy    int
	LongPassing        int
	LongShots          int
	Penalties          int
	ShortPassing       int
	ShotPower          int
	SlidingTackle      int
	StandingTackle     int
	Volleys            int
}

type GoalkeepingAttributes struct {
	ID int `json:"-"`

	GKDiving      int
	GKHandling    int
	GKKicking     int
	GKPositioning int
	GKReflexes    int
}
