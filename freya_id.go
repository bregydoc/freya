package main

type FreyaIDGenerator struct {
	alphabet string
}

func (id *FreyaIDGenerator) GetNewFreyaID() (string, error) {
	return Generate(id.alphabet, 8)
}

func (id *FreyaIDGenerator) GetLittleHash() (string, error) {
	return Generate(id.alphabet, 4)
}
