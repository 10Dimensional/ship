package integration

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/docker/docker/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/replicatedhq/ship/pkg/cli"
)

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "integration")
}

// CompareDir returns false if the two directories have different contents
func CompareDir(expected, actual string) (bool, error) {
	expectedDir, err := ioutil.ReadDir(expected)
	Expect(err).NotTo(HaveOccurred())

	actualDir, err := ioutil.ReadDir(actual)
	Expect(err).NotTo(HaveOccurred())

	Expect(len(actualDir)).To(Equal(len(expectedDir)), fmt.Sprintf("Number of files in %s and %s differed", expected, actual))

	expectedMap := make(map[string]os.FileInfo)
	for _, file := range expectedDir {
		expectedMap[file.Name()] = file
	}

	actualMap := make(map[string]os.FileInfo)
	for _, file := range expectedDir {
		actualMap[file.Name()] = file
	}

	for name, expectedFile := range expectedMap {
		actualFile, ok := actualMap[name]
		Expect(ok).To(BeTrue())
		Expect(actualFile.IsDir()).To(Equal(expectedFile.IsDir()))

		expectedFilePath := filepath.Join(expected, expectedFile.Name())
		actualFilePath := filepath.Join(actual, actualFile.Name())

		if expectedFile.IsDir() {
			// compare child items
			result, err := CompareDir(expectedFilePath, actualFilePath)
			if !result || err != nil {
				return result, err
			}
		} else {
			// compare expectedFile contents
			expectedContents, err := ioutil.ReadFile(expectedFilePath)
			Expect(err).NotTo(HaveOccurred())
			actualContents, err := ioutil.ReadFile(actualFilePath)
			Expect(err).NotTo(HaveOccurred())

			Expect(string(actualContents)).To(Equal(string(expectedContents)), fmt.Sprintf("Contents of files %s (expected) and %s (actual) did not match", expectedFilePath, actualFilePath))
		}
	}

	return true, nil
}

var _ = Describe("basic", func() {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	dockerClient.NegotiateAPIVersion(context.Background())

	integrationDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	files, err := ioutil.ReadDir(integrationDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			Context(fmt.Sprintf("When the spec in %s is run", file.Name()), func() {
				testPath := path.Join(integrationDir, file.Name())
				testOutputPath := path.Join(testPath, "tmp")
				testInputPath := path.Join(testPath, "input")

				custIdBytes, err := ioutil.ReadFile(path.Join(testPath, "customer_id"))
				if err != nil {
					panic(err)
				}
				customerId := string(custIdBytes)

				installationIdBytes, err := ioutil.ReadFile(path.Join(testPath, "installation_id"))
				if err != nil {
					panic(err)
				}
				installationId := string(installationIdBytes)

				channelIdBytes, err := ioutil.ReadFile(path.Join(testPath, "channel_id"))
				if err != nil {
					panic(err)
				}
				channelId := string(channelIdBytes)

				releaseVersionBytes, err := ioutil.ReadFile(path.Join(testPath, "release_version"))
				if err != nil {
					panic(err)
				}
				releaseVersion := string(releaseVersionBytes)

				BeforeEach(func() {
					//create a temporary directory within this directory to compare files with
					os.RemoveAll(testOutputPath)
					err := os.Mkdir(testOutputPath, os.ModeDir|os.ModePerm)
					if err != nil {
						panic(err)
					}
					os.Chdir(testOutputPath)
				})

				AfterEach(func() {
					//remove the temporary directory
					err := os.RemoveAll(testOutputPath)
					if err != nil {
						panic(err)
					}
					os.Chdir(integrationDir)
				})

				It("Should output files matching those expected when running in local mode", func() {
					cmd := cli.RootCmd()
					buf := new(bytes.Buffer)
					cmd.SetOutput(buf)
					cmd.SetArgs(append([]string{
						"--headless",
						fmt.Sprintf("--studio-file=%s", path.Join(testInputPath, ".ship/release.yml")),
						fmt.Sprintf("--state-file=%s", path.Join(testInputPath, ".ship/state.json")),
						"--log-level=off",
					}))
					err := cmd.Execute()
					Expect(err).NotTo(HaveOccurred())

					//compare the files in the temporary directory with those in the "expected" directory
					result, err := CompareDir(path.Join(testPath, "expected"), testOutputPath)
					Expect(err).NotTo(HaveOccurred())
					Expect(result).To(BeTrue())

				}, 60)

				It("Should output files matching those expected when communicating with the graphql api", func() {
					cmd := cli.RootCmd()
					buf := new(bytes.Buffer)
					cmd.SetOutput(buf)
					cmd.SetArgs(append([]string{
						"--headless",
						fmt.Sprintf("--state-file=%s", path.Join(testInputPath, ".ship/state.json")),
						"--customer-endpoint=https://pg.staging.replicated.com/graphql",
						"--log-level=off",
						fmt.Sprintf("--customer-id=%s", customerId),
						fmt.Sprintf("--installation-id=%s", installationId),
						fmt.Sprintf("--channel-id=%s", channelId),
						fmt.Sprintf("--release-semver=%s", releaseVersion),
					}))
					err := cmd.Execute()
					Expect(err).NotTo(HaveOccurred())

					//compare the files in the temporary directory with those in the "expected" directory
					result, err := CompareDir(path.Join(testPath, "expected"), testOutputPath)
					Expect(err).NotTo(HaveOccurred())
					Expect(result).To(BeTrue())

				}, 60)
			})
		}
	}
})
