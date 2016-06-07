package main_test

import (
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Out", func() {
	var (
		pathToBuiltBinary string
		cmd               *exec.Cmd
	)

	BeforeSuite(func() {
		var err error
		pathToBuiltBinary, err = gexec.Build("github.com/nickwei84/sms-resource/out")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	BeforeEach(func() {
		cmd = exec.Command(pathToBuiltBinary)
	})

	Context("when stdin input is invalid", func() {
		Context("because it is empty", func() {
			It("should output an error to stderr", func() {
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session).Should(gexec.Exit(1))
				Eventually(session.Err).Should(gbytes.Say("error parsing stdin as JSON: unexpected end of JSON input"))
				Eventually(session.Out).Should(gbytes.Say(""))
			})
		})

		Context("because it is invalid JSON", func() {
			It("should output an error to stderr", func() {
				cmd.Stdin = strings.NewReader(`
{
	malformed: "JSON"
}
`)
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session).Should(gexec.Exit(1))
				Eventually(session.Err).Should(gbytes.Say("error parsing stdin as JSON: invalid character 'm' looking for beginning of object key string"))
				Eventually(session.Out).Should(gbytes.Say(""))
			})
		})

		Context("because field values are invalid", func() {
			It("should output an error to stderr", func() {
				cmd.Stdin = strings.NewReader(`
{
	"source": {
		"aws_access_key_id": "",
		"aws_secret_access_key": "key123",
		"topic": "concourse"
	},
	"params": {
		"subscribers": [
			"1234567890"
		],
		"message": "hello!"
	}
}
`)
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
				Eventually(session).Should(gexec.Exit(1))
				Eventually(session.Err).Should(gbytes.Say("source.aws_access_key_id from stdin is either empty or missing"))
				Eventually(session.Out).Should(gbytes.Say(""))
			})
		})
	})
})
