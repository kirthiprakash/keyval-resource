package main_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"gstack.io/concourse/keyval-resource/models"
)

var _ = Describe("In", func() {
	var (
		tmpdir      string
		destination string
		inCmd       *exec.Cmd
	)

	BeforeEach(func() {
		var err error

		tmpdir, err = ioutil.TempDir("", "in-destination")
		Expect(err).NotTo(HaveOccurred())

		destination = path.Join(tmpdir, "in-dir")

		inCmd = exec.Command(inPath, destination)
	})

	AfterEach(func() {
		os.RemoveAll(tmpdir)
	})

	Context("when executed", func() {
		var (
			request  models.InRequest
			response models.InResponse
		)

		BeforeEach(func() {

			request = models.InRequest{
				Version: models.Version{
					"a": "1",
					"b": "2",
				},
				Source: models.Source{},
			}

			response = models.InResponse{}
		})

		JustBeforeEach(func() {
			stdin, err := inCmd.StdinPipe()
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(inCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			err = json.NewEncoder(stdin).Encode(request)
			Expect(err).NotTo(HaveOccurred())

			<-session.Exited
			Expect(session.ExitCode()).To(Equal(0))

			err = json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).NotTo(HaveOccurred())
		})

		It("reports the version to be the same as the version in input", func() {
			Expect(len(response.Version)).To(Equal(2))
			Expect(response.Version["a"]).To(Equal("1"))
			Expect(response.Version["b"]).To(Equal("2"))
		})

		It("writes key-value pairs to files in the destination directory", func() {
			aBytes, err := ioutil.ReadFile(filepath.Join(destination, "a"))
			Expect(err).NotTo(HaveOccurred())
			bBytes, err := ioutil.ReadFile(filepath.Join(destination, "b"))
			Expect(err).NotTo(HaveOccurred())

			Expect(string(aBytes)).To(Equal("1"))
			Expect(string(bBytes)).To(Equal("2"))
		})

		Context("when the request has no keys in version", func() {
			BeforeEach(func() {
				request.Version = models.Version{}
			})

			It("reports empty data", func() {
				Expect(len(response.Version)).To(Equal(0))

				files, err := ioutil.ReadDir(destination)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(files)).To(Equal(0))
			})
		})
	})
})
