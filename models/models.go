package models

type Reserva struct {
	Vuelos    []Vuelo    `json:"vuelos,omitempty" bson:"vuelos,omitempty"`
	Pasajeros []Pasajero `json:"pasajeros,omitempty" bson:"pasajeros,omitempty"`
	PNR       string     `json:"pnr,omitempty" bson:"pnr,omitempty"`
}

type Pasajero struct {
	Nombre      string                `json:"nombre,omitempty" bson:"nombre,omitempty"`
	Apellido    string                `json:"apellido,omitempty" bson:"apellido,omitempty"`
	Edad        int                   `json:"edad,omitempty" bson:"edad,omitempty"`
	Ancillaries []Ancillarie_pasajero `json:"ancillaries,omitempty" bson:"ancillaries,omitempty"`
	Balances    Balance               `json:"balances,omitempty" bson:"balances,omitempty"`
}

type Ancillarie_pasajero struct {
	Ida []struct {
		Ssr      string
		Cantidad int
	} `json:"ida,omitempty" bson:"ida,omitempty"`
	Vuelta []struct {
		Ssr      string
		Cantidad int
	} `json:"vuelta,omitempty" bson:"vuelta,omitempty"`
}

type Balance struct {
	AncillariesIda    int `json:"ancillariesida,omitempty" bson:"ancillariesida,omitempty"`
	VueloIda          int `json:"vueloida,omitempty" bson:"vueloida,omitempty"`
	AncillariesVuelta int `json:"ancillariesvuelta,omitempty" bson:"ancillariesvuelta,omitempty"`
	VueloVuelta       int `json:"vuelovuelta,omitempty" bson:"vuelovuelta,omitempty"`
}

type Vuelo struct {
	NumeroVuelo string       `json:"numero_vuelo,omitempty" bson:"numero_vuelo,omitempty"`
	Origen      string       `json:"origen,omitempty" bson:"origen,omitempty"`
	Destino     string       `json:"destino,omitempty" bson:"destino,omitempty"`
	HoraSalida  string       `json:"hora_salida,omitempty" bson:"hora_salida,omitempty"`
	HoraLlegada string       `json:"hora_llegada,omitempty" bson:"hora_llegada,omitempty"`
	Fecha       string       `json:"fecha,omitempty" bson:"fecha,omitempty"`
	Avion       Avion        `json:"avion,omitempty" bson:"avion,omitempty"`
	Ancillaries []Ancillarie `json:"ancillaries,omitempty" bson:"ancillaries,omitempty"`
}

type Ancillarie struct {
	Nombre string `json:"nombre,omitempty" bson:"nombre,omitempty"`
	Stock  int    `json:"stock,omitempty" bson:"stock,omitempty"`
	Ssr    string `json:"ssr,omitempty" bson:"ssr,omitempty"`
}

type Avion struct {
	Modelo           string `json:"modelo,omitempty" bson:"modelo,omitempty"`
	NumeroDeSerie    string `json:"numero_de_serie,omitempty" bson:"numero_de_serie,omitempty"`
	StockDePasajeros int    `json:"stock_de_pasajeros,omitempty" bson:"stock_de_pasajeros,omitempty"`
}

type PNR struct {
	Value string `json:"value,omitempty" bson:"value,omitempty"`
}
