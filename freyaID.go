package main

const freyaAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GetNewFreyaID() (string, error) {
	return Generate(freyaAlphabet, 8)
}

func GetLittleHash() (string, error) {
	return Generate(freyaAlphabet, 4)
}
