package models

type Reserva struct {
	Vuelos    []Vuelo    `json:"vuelos,omitempty" bson:"vuelos,omitempty"`
	Pasajeros []Pasajero `json:"pasajeros,omitempty" bson:"pasajeros,omitempty"`
	PNR       string     `json:"pnr,omitempty" bson:"pnr,omitempty"`
}

type Pasajero struct {
	Nombre      string               `json:"nombre,omitempty" bson:"nombre,omitempty"`
	Apellido    string               `json:"apellido,omitempty" bson:"apellido,omitempty"`
	Edad        int                  `json:"edad,omitempty" bson:"edad,omitempty"`
	Ancillaries []AncillariePasajero `json:"ancillaries,omitempty" bson:"ancillaries,omitempty"`
	Balances    Balance              `json:"balances,omitempty" bson:"balances,omitempty"`
}

type AncillariePasajero struct {
	Ida    []AncillarieDetail `json:"ida,omitempty" bson:"ida,omitempty"`
	Vuelta []AncillarieDetail `json:"vuelta,omitempty" bson:"vuelta,omitempty"`
}

type AncillarieDetail struct {
	Ssr      string `json:"ssr,omitempty" bson:"ssr,omitempty"`
	Cantidad int    `json:"cantidad,omitempty" bson:"cantidad,omitempty"`
}

type Balance struct {
	AncillariesIda    int `json:"ancillaries_ida,omitempty" bson:"ancillaries_ida,omitempty"`
	VueloIda          int `json:"vuelo_ida,omitempty" bson:"vuelo_ida,omitempty"`
	AncillariesVuelta int `json:"ancillaries_vuelta,omitempty" bson:"ancillaries_vuelta,omitempty"`
	VueloVuelta       int `json:"vuelo_vuelta,omitempty" bson:"vuelo_vuelta,omitempty"`
}

type Vuelo struct {
	NumeroVuelo string       `json:"numero_vuelo,omitempty" bson:"numero_vuelo,omitempty"`
	Origen      string       `json:"origen,omitempty" bson:"origen,omitempty"`
	Destino     string       `json:"destino,omitempty" bson:"destino,omitempty"`
	HoraSalida  string       `json:"hora_salida,omitempty" bson:"hora_salida,omitempty"`
	HoraLlegada string       `json:"hora_llegada,omitempty" bson:"hora_llegada,omitempty"`
	Fecha       string       `json:"fecha,omitempty" bson:"fecha,omitempty"`
	Avion       *Avion       `json:"avion,omitempty" bson:"avion,omitempty"`
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
	Pnr string `json:"pnr,omitempty" bson:"pnr,omitempty"`
}

type AncillarieData struct {
	Nombre string `json:"nombre,omitempty" bson:"nombre,omitempty"`
	Precio int    `json:"precio,omitempty" bson:"precio,omitempty"`
	Ssr    string `json:"ssr,omitempty" bson:"ssr,omitempty"`
}

type Stock struct {
	Stock int `json:"stock_de_pasajeros" bson:"stock_de_pasajeros"`
}

type Compra struct {
	Origen            string `json:"ida,omitempty" bson:"ida,omitempty"`
	Destino           string `json:"vuelta,omitempty" bson:"vuelta,omitempty"`
	TotalPasajes      int    `json:"total_pasajes,omitempty" bson:"total_pasajes,omitempty"`
	CantidadPasajeros int    `json:"cantidad_pasajeros,omitempty" bson:"cantidad_pasajeros,omitempty"`
	FechaVuelo        string `json:"fecha_vuelo,omitempty" bson:"fecha_vuelo,omitempty"`
}

type Statistics struct {
	RutaMayorGanancia  string                `json:"ruta_mayor_ganancia,omitempty" bson:"ruta_mayor_ganancia,omitempty"`
	RutaMenorGanancia  string                `json:"ruta_menor_ganancia,omitempty" bson:"ruta_menor_ganancia,omitempty"`
	RankingAncillaries []StatisticAncillarie `json:"ranking_ancillaries,omitempty" bson:"ranking_ancillaries,omitempty"`
	PromedioPasajeros  PromedioPasajeros     `json:"promedio_pasajeros" bson:"promedio_pasajeros"`
}

type StatisticAncillarie struct {
	Nombre   string `json:"nombre,omitempty" bson:"nombre,omitempty"`
	Ganancia int    `json:"ganancia,omitempty" bson:"ganancia,omitempty"`
	Ssr      string `json:"ssr,omitempty" bson:"ssr,omitempty"`
}

type PromedioPasajeros struct {
	Enero      int `json:"enero" bson:"enero"`
	Febrero    int `json:"febrero" bson:"febrero"`
	Marzo      int `json:"marzo" bson:"marzo"`
	Abril      int `json:"abril" bson:"abril"`
	Mayo       int `json:"mayo" bson:"mayo"`
	Junio      int `json:"junio" bson:"junio"`
	Julio      int `json:"julio" bson:"julio"`
	Agosto     int `json:"agosto" bson:"agosto"`
	Septiembre int `json:"septiembre" bson:"septiembre"`
	Octubre    int `json:"octubre" bson:"octubre"`
	Noviembre  int `json:"noviembre" bson:"noviembre"`
	Diciembre  int `json:"diciembre" bson:"diciembre"`
}
