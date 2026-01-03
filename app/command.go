package main

type Commander interface {
	Execute()
	GetType() string
}
