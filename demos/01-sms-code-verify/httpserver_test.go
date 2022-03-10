package main

import (
	"testing"

	. "github.com/onsi/gomega"
)

func Test_GenSmsCode(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := genSmsCode()
		NewWithT(t).Expect(len(s)).To(Equal(6))
	}
}
