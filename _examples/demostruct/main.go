package main

import (
	"fmt"
	"log/slog"
)

type RepoA struct{}
type RepoB struct{}
type RepoC struct{}

type ServiceA struct {
	a *RepoA
	b *RepoB
}

func NewServiceA() *ServiceA {
	return &ServiceA{}
}

func (s *ServiceA) UseRepo(r interface{}) *ServiceA {
	switch repo := r.(type) {
	case *RepoA:
		s.a = repo
	case *RepoB:
		s.b = repo
	default:
		fmt.Println("does not support this repo")
	}
	return s
}

type ServiceB struct {
	a *RepoA
	b *RepoB
	c *RepoC
}

func NewServiceB(args ...interface{}) *ServiceB {
	s := &ServiceB{}
	for _, arg := range args {
		switch repo := arg.(type) {
		case *RepoA:
			s.a = repo
		case *RepoB:
			s.b = repo
		case *RepoC:
			s.c = repo
		default:
			fmt.Println("does not support this repo")
		}
	}
	return s
}

func main() {
	sA := NewServiceA()
	sA.UseRepo(&RepoA{}).
		UseRepo(&RepoB{})

	sB := NewServiceB(&RepoA{}, &RepoB{}, &RepoC{})
	slog.Info("success", slog.Any("serviceA", sA), slog.Any("serviceB", sB))
}
