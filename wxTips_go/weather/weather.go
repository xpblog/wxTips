package weather

type WeatherRes struct {
	Data Data `json:"data"`
}
type Data struct {
	Forecast []Forecast `json:"forecast"`
}

type Forecast struct {
	High string `json:"high"`
	Low  string `json:"low"`
	Type string `json:"type"`
	Fl   string `json:"fl"`
}
