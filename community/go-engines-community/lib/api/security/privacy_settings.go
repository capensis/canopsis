package security

type PrivacySettings struct {
	Author    string `bson:"author"`
	IsPrivate bool   `bson:"is_private"`
}

type ViewTabPrivacySettings struct {
	PrivacySettings `bson:"inline"`
	View            string `bson:"view"`
}
