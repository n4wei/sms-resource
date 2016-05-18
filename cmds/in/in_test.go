package main_test

import (
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("In", func() {
	var pathToBuiltBinary string

	BeforeSuite(func() {
		var err error
		pathToBuiltBinary, err = gexec.Build("github.com/nickwei84/sms-resource/cmds/in")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	Context("when none or invalid input is passed in via stdin", func() {
		var cmd *exec.Cmd
		var session *gexec.Session

		BeforeEach(func() {
			var err error
			cmd = exec.Command(pathToBuiltBinary)
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should output an error to stderr", func() {
			Eventually(session).Should(gexec.Exit(1))
			Eventually(session.Err).Should(gbytes.Say("Error unmarshalling JSON: unexpected end of JSON input"))
			Eventually(session.Out).Should(gbytes.Say(""))
		})
	})

	Context("when version is not passed in via stdin", func() {
		var cmd *exec.Cmd
		var session *gexec.Session

		BeforeEach(func() {
			var err error
			cmd = exec.Command(pathToBuiltBinary)
			cmd.Stdin = strings.NewReader(`
{
	"some-key": "some-value"
}
`)
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should output an error to stderr", func() {
			Eventually(session).Should(gexec.Exit(1))
			Eventually(session.Err).Should(gbytes.Say("Error: version key pair is missing from stdin"))
			Eventually(session.Out).Should(gbytes.Say(""))
		})
	})

	Context("when version is passed in via stdin", func() {
		var cmd *exec.Cmd
		var session *gexec.Session

		BeforeEach(func() {
			var err error
			cmd = exec.Command(pathToBuiltBinary)
			cmd.Stdin = strings.NewReader(`
{
	"version": "abc123"
}
`)
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should output version to stdout", func() {
			Eventually(session).Should(gexec.Exit(0))
			Eventually(session.Out).Should(gbytes.Say(`{"version":"abc123"}`))
			Eventually(session.Err).Should(gbytes.Say(""))
		})
	})
})
