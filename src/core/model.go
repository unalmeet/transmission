package core

type Client struct {
	IdMeeting string `json:"idMeeting" bson:"idMeeting" validate:"empty=false"`
	IdUser    int    `json:"idUser" bson:"idUser"`
	IdSession int    `json:"idSession" bson:"idSession"`
	Token     string `json:"token" bson:"token"`
	Media     []byte `json:"media" bson:"media"`
}

type Broadcast struct {
	Token     string `json:"token" bson:"token" validate:"empty=false"`
	IdSession []int  `json:"idSession" bson:"idSession"`
	Media     []byte `json:"media" bson:"media"`
}