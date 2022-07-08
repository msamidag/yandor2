package model

type Calisan struct {
	ID          int       `json:"id" bson:"id"`
	NameSurname string    `json:"nameSurname" bson:"nameSurname"`
	Soygecmis   Soygecmis `json:"soygecmis" bson:"soygecmis"`
	Ozgecmis    Ozgecmis  `json:"ozgecmis" bson:"ozgecmis"`
	Anamnez     Anamnez   `json:"anamnez" bson:"anamnez"`
}

type Soygecmis struct {
	Calisanid int    `json:"calisanid" bson:"calisanid"`
	Anne      string `json:"anne" bson:"anne"`
	Baba      string `json:"baba" bson:"baba"`
}

type Ozgecmis struct {
	CalisanID  int    `json:"calisanid" bson:"calisanid"`
	KrHastalik string `json:"krHastalik" bson:"krHastalik"`
	Tetanoz    string `json:"tetanoz" bson:"tetanoz"`
}
type Anamnez struct {
	CalisanId int    `json:"calisanid" bson:"calisanid"`
	Sigara    string `json:"sigara" bson:"sigara"`
	Alkol     string `json:"alkol" bson:"alkol"`
	Ameliyat  string `json:"ameliyat" bson:"ameliyat"`
}
