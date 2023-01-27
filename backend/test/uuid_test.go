package test

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestGeneratorUUID(t *testing.T) {
	s := uuid.NewV4().String()
	fmt.Println(s)
	//f5f2689d-7872-4046-9c09-b857d238ab2e
	//b1f51ff9-4edc-4ca1-8a4a-9367e63d49a7
	//len: 36
}
