package tests

import (
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/portworx/torpedo/drivers/node"
	"github.com/portworx/torpedo/drivers/scheduler"
	. "github.com/portworx/torpedo/tests"
	"github.com/sirupsen/logrus"
)

const (
	defaultTimeout       = 5 * time.Minute
	defaultRetryInterval = 20 * time.Second
	appReadinessTimeout  = 20 * time.Minute
)

func TestPxcentral(t *testing.T) {
	RegisterFailHandler(Fail)

	var specReporters []Reporter
	junitReporter := reporters.NewJUnitReporter("/testresults/junit_basic.xml")
	specReporters = append(specReporters, junitReporter)
	RunSpecsWithDefaultAndCustomReporters(t, "Torpedo : px-central", specReporters)
}

var _ = BeforeSuite(func() {
	logrus.Infof("Init instance")
	InitInstance()
})

// This test performs basic test of installing px-central with helm
// px-license-server and px-minotor will be installed after px-central is validated
var _ = Describe("{Installpxcentral}", func() {
	It("has to setup, validate and teardown apps", func() {
		var context *scheduler.Context
		lsOptions := scheduler.ScheduleOptions{
			AppKeys:            []string{"px-license-server"},
			StorageProvisioner: Inst().Provisioner,
		}

		Step("Install px-central using the px-backup helm chart then validate", func() {
			contexts, err := Inst().S.Schedule(Inst().InstanceID, scheduler.ScheduleOptions{
				AppKeys:            []string{"px-central"},
				StorageProvisioner: Inst().Provisioner,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(contexts).NotTo(BeEmpty())

			// Skipping volume validation until other volume providers are implemented.
			// Also change the app readinessTimeout to 20mins
			context = contexts[0]
			context.SkipVolumeValidation = true
			context.ReadinessTimeout = appReadinessTimeout

			ValidateContext(context)
			logrus.Infof("Successfully validated specs for px-central")
		})

		Step("Install px-license-server then validate", func() {
			// label px/ls=true on 2 worker nodes
			for i, node := range node.GetWorkerNodes() {
				if i == 2 {
					break
				}
				err := Inst().S.AddLabelOnNode(node, "px/ls", "true")
				Expect(err).NotTo(HaveOccurred())
			}

			err := Inst().S.AddTasks(context, lsOptions)
			Expect(err).NotTo(HaveOccurred())

			ValidateContext(context)
			logrus.Infof("Successfully validated specs for px-license-server")
		})

		Step("Uninstall license server and monitoring", func() {
			err := Inst().S.ScheduleUninstall(context, lsOptions)
			Expect(err).NotTo(HaveOccurred())

			ValidateContext(context)
			logrus.Infof("Successfully uninstalled px-license-server")
		})

		Step("destroy apps", func() {
			opts := make(map[string]bool)
			opts[scheduler.OptionsWaitForResourceLeakCleanup] = true

			TearDownContext(context, opts)
			logrus.Infof("Successfully destroyed px-central")
		})
	})
})

var _ = AfterSuite(func() {
	PerformSystemCheck()
	ValidateCleanup()
})

func TestMain(m *testing.M) {
	ParseFlags()
	os.Exit(m.Run())
}
