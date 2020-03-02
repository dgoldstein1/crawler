package db

type GraphResponseError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type GraphResponseSuccess struct {
	NeighborsAdded []string `json:"neighborsAdded"`
}

type TwoWayEntry struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type TwoWayResponse struct {
	Errors  []string      `json:"errors"`
	Entries []TwoWayEntry `json:"entries"`
}
