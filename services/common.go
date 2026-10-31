package services

type Provider func(string) (Quote, error)
