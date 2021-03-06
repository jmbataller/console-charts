package installer

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/lightbend/console-charts/enterprise-suite/gotests/args"
	"github.com/lightbend/console-charts/enterprise-suite/gotests/util/lbc"

	"github.com/lightbend/console-charts/enterprise-suite/gotests/testenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInstaller(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Installer (lbc.py) Suite")
}

var _ = BeforeSuite(func() {
	testenv.InitEnv()
})

var _ = AfterSuite(func() {
	testenv.CloseEnv()
})

func write(file *os.File, content string) {
	if _, err := file.Write([]byte(content)); err != nil {
		panic(err)
	}
}

var _ = Describe("all:lbc.py", func() {
	var (
		valuesFile *os.File
	)

	BeforeEach(func() {
		var err error
		valuesFile, err = ioutil.TempFile("", "values-*.yaml")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		err := os.Remove(valuesFile.Name())
		Expect(err).To(Succeed())
	})

	Context("upgrades", func() {
		Context("disable persistent volumes", func() {
			var installer *lbc.Installer

			BeforeEach(func() {
				preInstaller := lbc.DefaultInstaller()
				preInstaller.UsePersistentVolumes = "true"
				Expect(preInstaller.Install()).To(Succeed(), "install with PVs")

				write(valuesFile, `usePersistentVolumes: false`)
				installer = lbc.DefaultInstaller()
				installer.AdditionalHelmArgs = []string{"-f " + valuesFile.Name()}
			})

			It("should fail if we don't provide --delete-pvcs", func() {
				installer.UsePersistentVolumes = ""
				installer.ForceDeletePVCs = false
				Expect(installer.Install()).ToNot(Succeed())
			})

			It("should succeed if we provide --delete-pvcs", func() {
				installer.UsePersistentVolumes = ""
				installer.ForceDeletePVCs = true
				Expect(installer.Install()).To(Succeed())
			})
		})
	})

	Context("arg parsing", func() {
		It("should fail if conflicting namespaces", func() {
			installer := lbc.DefaultInstaller()
			installer.AdditionalLBCArgs = []string{"--namespace=" + args.ConsoleNamespace}
			installer.AdditionalHelmArgs = []string{"--namespace=my-busted-namespace"}
			Expect(installer.Install()).ToNot(Succeed())
		})
	})
})
