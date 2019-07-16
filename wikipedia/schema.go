package wikipedia


type GraphResponseError struct {
	Code int `json:"code"`
	Error string `json:"error"`
}

type GraphResponseSuccess struct {
	NeighborsAdded []int `json:"neighborsAdded"`
}

type PropertiesResponse struct {
	Parse PropertiesValues `json:"parse"`
}
type PropertiesValues struct {
	Pageid int `json:"pageid"`
	// ...drop title and properties keys
}
