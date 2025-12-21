package main

type ICommand interface {
	Execute()
	GetType() string
}
