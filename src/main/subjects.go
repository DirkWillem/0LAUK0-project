package main

import "main/dispatch"

var (
	dispatcher *dispatch.Dispatcher
	medicationsSubject *dispatch.CollectionSubject
)

func init() {
	dispatcher = dispatch.NewDispatcher()

	medicationsSubject = dispatch.NewCollectionSubject("medications", dispatcher)
}
