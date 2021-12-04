package core

type Client struct {
	IdMeeting string `json:"idMeeting" bson:"idMeeting" validate:"empty=false"`
	IdSession int    `json:"idSession" bson:"idSession"`
	Media     []byte `json:"media" bson:"media"`
}

type Broadcast struct {
	IdMeeting string `json:"idMeeting" bson:"idMeeting" validate:"empty=false"`
	IdSession []int  `json:"idSession" bson:"idSession"`
	Media     []byte `json:"media" bson:"media"`
}