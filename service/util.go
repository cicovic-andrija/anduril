package service

type TraceCallback func(string, ...interface{})

type Task func(TraceCallback, ...interface{}) error
