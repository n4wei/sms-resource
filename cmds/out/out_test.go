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
	var pathToBuiltBinary string

	BeforeSuite(func() {
		var err error
		pathToBuiltBinary, err = gexec.Build("github.com/nickwei84/sms-resource/cmds/out")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	Context("when no input is passed in via stdin", func() {
		It("should output an error to stderr", func() {
			cmd := exec.Command(pathToBuiltBinary)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
			Eventually(session.Err).Should(gbytes.Say("Error parsing stdin as JSON: unexpected end of JSON input"))
			Eventually(session.Out).Should(gbytes.Say(""))
		})
	})

	Context("when invalid JSON input is passed in via stdin", func() {
		It("should output an error to stderr", func() {
			cmd := exec.Command(pathToBuiltBinary)
			cmd.Stdin = strings.NewReader(`
{
	malformed: "JSON"
}
`)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
			Eventually(session.Err).Should(gbytes.Say("Error parsing stdin as JSON: invalid character 'm' looking for beginning of object key string"))
			Eventually(session.Out).Should(gbytes.Say(""))
		})
	})

	Context("when inputs from stdin are invalid", func() {
		var cmd *exec.Cmd

		BeforeEach(func() {
			cmd = exec.Command(pathToBuiltBinary)
		})

		It("should output an error to stderr if AWS key ID is missing", func() {
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
			Eventually(session.Err).Should(gbytes.Say("Error: source.aws_access_key_id from stdin is either empty or missing"))
			Eventually(session.Out).Should(gbytes.Say(""))
		})

		It("should output an error to stderr if AWS secret key is missing", func() {
			cmd.Stdin = strings.NewReader(`
{
	"source": {
		"aws_access_key_id": "id123",
		"aws_secret_access_key": "",
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
			Eventually(session.Err).Should(gbytes.Say("Error: source.aws_secret_access_key from stdin is either empty or missing"))
			Eventually(session.Out).Should(gbytes.Say(""))
		})

		It("should output an error to stderr if topic is missing", func() {
			cmd.Stdin = strings.NewReader(`
{
	"source": {
		"aws_access_key_id": "id123",
		"aws_secret_access_key": "key123",
		"topic": ""
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
			Eventually(session.Err).Should(gbytes.Say("Error: source.topic from stdin is either empty or missing"))
			Eventually(session.Out).Should(gbytes.Say(""))
		})

		It("should output an error to stderr if topic exceeds max character limit", func() {
			cmd.Stdin = strings.NewReader(`
{
	"source": {
		"aws_access_key_id": "id123",
		"aws_secret_access_key": "key123",
		"topic": "concourse1234567890"
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
			Eventually(session.Err).Should(gbytes.Say("Error: source.topic from stdin cannot exceed 10 characters"))
			Eventually(session.Out).Should(gbytes.Say(""))
		})

		It("should output an error to stderr if no subscribers are provided", func() {
			cmd.Stdin = strings.NewReader(`
{
	"source": {
		"aws_access_key_id": "id123",
		"aws_secret_access_key": "key123",
		"topic": "concourse"
	},
	"params": {
		"subscribers": [],
		"message": "hello!"
	}
}
`)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
			Eventually(session.Err).Should(gbytes.Say("Error: params.subscribers from stdin is either empty or missing"))
			Eventually(session.Out).Should(gbytes.Say(""))
		})

		It("should output an error to stderr if message is missing", func() {
			cmd.Stdin = strings.NewReader(`
{
	"source": {
		"aws_access_key_id": "id123",
		"aws_secret_access_key": "key123",
		"topic": "concourse"
	},
	"params": {
		"subscribers": [
			"1234567890"
		],
		"message": ""
	}
}
`)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(1))
			Eventually(session.Err).Should(gbytes.Say("Error: params.message from stdin is either empty or missing"))
			Eventually(session.Out).Should(gbytes.Say(""))
		})
	})

	Context("when input to stdin is valid", func() {
	})
})
