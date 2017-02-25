package main

import "main/dispatch"

var (
	dispatcher *dispatch.Dispatcher
	medicationsSubject *dispatch.CollectionSubject
	dosesSubject *DosesSubject
)

func init() {
	dispatcher = dispatch.NewDispatcher()

	medicationsSubject = dispatch.NewCollectionSubject("medications", dispatcher)
	dosesSubject = NewDosesSubject(dispatcher)
}
