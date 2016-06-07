package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Check", func() {
	var pathToBuiltBinary string

	BeforeSuite(func() {
		var err error
		pathToBuiltBinary, err = gexec.Build("github.com/nickwei84/sms-resource/check")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	It("should output an empty JSON list to stdout", func() {
		cmd := exec.Command(pathToBuiltBinary)
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Eventually(session.Out).Should(gbytes.Say(`\[\]`))
		Eventually(session.Err).Should(gbytes.Say(""))
	})
})
