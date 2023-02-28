package main

type register interface {
	RegisterByJson(Person, string) error
}
