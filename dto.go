package main

type Sub struct {
	Action  string `json:"action"`
	Symbols string `json:"symbols"`
}

type Info struct {
	Timestamp string `json:"timestamp"`
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
}

type Command struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
}
