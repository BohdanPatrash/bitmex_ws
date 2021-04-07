package main

type Sub struct {
	Action  string   `json:"action"`
	Symbols []string `json:"symbols"`
}

type Info struct {
	Timestamp string  `json:"timestamp"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
}

type BitmexCommand struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
}

type Error struct {
	Message string `json:"error"`
}
