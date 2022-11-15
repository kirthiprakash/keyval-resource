package main_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"gstack.io/concourse/keyval-resource/models"
)

func writeFile(name string, contents string) {
	file, err := os.Create(name)

	Expect(err).NotTo(HaveOccurred())
	defer file.Close()
	writer := bufio.NewWriter(file)
	fmt.Fprint(writer, contents)
	writer.Flush()
}

var _ = Describe("Out", func() {
	var (
		sourceDir string
		outDir    = "out-dir"
		outCmd    *exec.Cmd
	)

	BeforeEach(func() {
		var err error

		sourceDir, err = ioutil.TempDir("", "out-source")
		Expect(err).NotTo(HaveOccurred())

		outCmd = exec.Command(outPath, sourceDir)

		os.MkdirAll(path.Join(sourceDir, outDir), 0755)
	})

	AfterEach(func() {
		os.RemoveAll(sourceDir)
	})

	Context("when executed", func() {
		var (
			request  models.OutRequest
			response models.OutResponse
		)

		BeforeEach(func() {
			request = models.OutRequest{}
			response = models.OutResponse{}
		})

		JustBeforeEach(func() {
			stdin := new(bytes.Buffer)

			err := json.NewEncoder(stdin).Encode(request)
			Expect(err).NotTo(HaveOccurred())

			outCmd.Stdin = stdin

			session, err := gexec.Start(outCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			<-session.Exited
			Expect(session.ExitCode()).To(Equal(0))

			err = json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("two files in artifact directory", func() {
			BeforeEach(func() {
				os.RemoveAll(outDir)
				writeFile(path.Join(sourceDir, outDir, "a"), "1")
				writeFile(path.Join(sourceDir, outDir, "b"), "2")

				request = models.OutRequest{
					Params: models.OutParams{
						Directory: outDir,
					},
				}
			})

			It("reports some value for UUID and UPDATED keys", func() {
				Expect(response.Version).To(HaveKey("UPDATED"))
				Expect(response.Version).To(HaveKey("UUID"))
				Expect(response.Version["UPDATED"]).NotTo(BeEmpty())
				Expect(response.Version["UUID"]).NotTo(BeEmpty())
			})

			It("reports key-value pairs from files into the version object", func() {
				Expect(len(response.Version)).To(Equal(4))
				Expect(response.Version["a"]).To(Equal("1"))
				Expect(response.Version["b"]).To(Equal("2"))
			})

			Context("when some value is overridden", func() {
				BeforeEach(func() {
					request.Params.Overrides = map[string]string{"a": "7"}
				})

				It("the values from 'put' step params overrides any value from files", func() {
					Expect(len(response.Version)).To(Equal(4))
					Expect(response.Version["a"]).To(Equal("7"))
					Expect(response.Version["b"]).To(Equal("2"))
				})
			})
		})

		Context("no files in artifact directory", func() {
			BeforeEach(func() {
				os.RemoveAll(outDir)

				request = models.OutRequest{
					Params: models.OutParams{
						Directory: outDir,
					},
				}
			})

			It("reports UUID and UPDATED with no key-value pair", func() {
				Expect(len(response.Version)).To(Equal(2))
				Expect(response.Version).To(HaveKey("UPDATED"))
				Expect(response.Version).To(HaveKey("UUID"))
				Expect(response.Version["UPDATED"]).To(Not(BeEmpty()))
				Expect(response.Version["UUID"]).To(Not(BeEmpty()))
			})
		})

	})

	Context("with invalid inputs", func() {
		var (
			request models.OutRequest
			session *gexec.Session
		)

		BeforeEach(func() {
			request = models.OutRequest{}
		})

		JustBeforeEach(func() {
			stdin := new(bytes.Buffer)

			err := json.NewEncoder(stdin).Encode(request)
			Expect(err).NotTo(HaveOccurred())

			outCmd.Stdin = stdin

			session, err = gexec.Start(outCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

		})

		Context("no file specified", func() {
			It("reports error", func() {
				<-session.Exited
				Expect(session.Err).To(gbytes.Say("missing parameter 'directory'"))
				Expect(session.ExitCode()).To(Equal(1))
			})
		})

	})
})
