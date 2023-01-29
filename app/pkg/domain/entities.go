package domain

type Location struct {
	Ident     int32 `json:"id"`
	Longitude int32 `json:"lon"`
	Altitude  int32 `json:"alt"`
	Timestamp int32 `json:"timestamp"`
}
