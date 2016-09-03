package main

type config struct {
	APIKey       string
	ListId       string
	DefaultStore string

	DataPath string
}

var Config = config{
	APIKey:       "e6e5da412f2a9676b00546d15d87dc5f-us4",
	ListId:       "23ad4e4ba4",
	DefaultStore: "7RtpEPYmCnJrnB",

	DataPath: "./data",
}
