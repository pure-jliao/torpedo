package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/portworx/torpedo/drivers/node"
	"github.com/portworx/torpedo/drivers/scheduler"
	. "github.com/portworx/torpedo/tests"

	// https://github.com/kubernetes/client-go/issues/242
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

const (
	scaleTimeout                = 10 * time.Minute
	autoNodeRecoveryTimeoutMins = 15 * time.Minute
)

func TestASG(t *testing.T) {
	RegisterFailHandler(Fail)

	var specReporters []Reporter
	junitReporter := reporters.NewJUnitReporter("/testresults/junit_basic.xml")
	specReporters = append(specReporters, junitReporter)
	RunSpecsWithDefaultAndCustomReporters(t, "Torpedo : ASG", specReporters)
}

var _ = BeforeSuite(func() {
	InitInstance()
})

// This test performs basic test of scaling up and down the asg cluster
var _ = Describe("{ClusterScaleUpDown}", func() {
	It("has to validate that storage nodes are not lost during asg scaledown", func() {

		var contexts []*scheduler.Context
		for i := 0; i < Inst().ScaleFactor; i++ {
			contexts = append(contexts, ScheduleAndValidate(fmt.Sprintf("asgscaleupdown-%d", i))...)
		}

		intitialNodeCount, err := Inst().N.GetASGClusterSize()
		Expect(err).NotTo(HaveOccurred())

		scaleupCount := intitialNodeCount + intitialNodeCount/2
		Step(fmt.Sprintf("scale up cluster from %d to %d nodes and validate",
			intitialNodeCount, (scaleupCount/3)*3), func() {

			Scale(scaleupCount)
			Step(fmt.Sprintf("validate number of storage nodes after scale up"), func() {
				ValidateClusterSize(scaleupCount)
			})

		})

		Step(fmt.Sprintf("scale down cluster back to original size of %d nodes",
			intitialNodeCount), func() {
			Scale(intitialNodeCount)

			Step(fmt.Sprintf("wait for %s minutes for auto recovery of storeage nodes",
				autoNodeRecoveryTimeoutMins.String()), func() {
				time.Sleep(autoNodeRecoveryTimeoutMins)
			})

			// After scale down, get fresh list of nodes
			// by re-initializing scheduler and volume driver
			err = Inst().S.RefreshNodeRegistry()
			Expect(err).NotTo(HaveOccurred())

			err = Inst().V.RefreshDriverEndpoints()
			Expect(err).NotTo(HaveOccurred())

			Step(fmt.Sprintf("validate number of storage nodes after scale down"), func() {
				ValidateClusterSize(intitialNodeCount)
			})
		})

		opts := make(map[string]bool)
		opts[scheduler.OptionsWaitForResourceLeakCleanup] = true
		ValidateAndDestroy(contexts, opts)

	})
})

var _ = AfterSuite(func() {
	PerformSystemCheck()
	CollectSupport()
	ValidateCleanup()
})

func init() {
	ParseFlags()
}

func Scale(count int64) {
	// In multi-zone ASG cluster, node count is per zone
	perZoneCount := count / 3

	err := Inst().N.SetASGClusterSize(perZoneCount, scaleTimeout)
	Expect(err).NotTo(HaveOccurred())
}

func ValidateClusterSize(count int64) {
	// In multi-zone ASG cluster, node count is per zone
	perZoneCount := count / 3

	// Validate total node count
	currentNodeCount, err := Inst().N.GetASGClusterSize()
	Expect(err).NotTo(HaveOccurred())
	Expect(currentNodeCount).Should(Equal(perZoneCount * 3))

	// Validate storage node count
	var expectedStorageNodesPerZone int
	if Inst().MaxStorageNodesPerAZ <= int(perZoneCount) {
		expectedStorageNodesPerZone = Inst().MaxStorageNodesPerAZ
	} else {
		expectedStorageNodesPerZone = int(perZoneCount)
	}
	storageNodes, err := getStorageNodes()
	Expect(err).NotTo(HaveOccurred())
	Expect(len(storageNodes)).Should(Equal(expectedStorageNodesPerZone * 3))
	logrus.Infof("Validated successfully that [%d] storage nodes are present", len(storageNodes))
}

func getStorageNodes() ([]node.Node, error) {

	storageNodes := []node.Node{}
	nodes := node.GetStorageDriverNodes()

	for _, node := range nodes {
		devices, err := Inst().V.GetStorageDevices(node)
		if err != nil {
			return nil, err
		}
		if len(devices) > 0 {
			storageNodes = append(storageNodes, node)
		}
	}
	return storageNodes, nil
}