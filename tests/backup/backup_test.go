package tests

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	api "github.com/portworx/px-backup-api/pkg/apis/v1"
	. "github.com/portworx/torpedo/tests"
)

const (
	orgID         = "simpleBackupOrg"
	BLocationName = "simpleBLocation"
	ClusterName   = "simpleBackupCluster"
	CredName      = "simpleBackupCred"
)

func TestBackup(t *testing.T) {
	RegisterFailHandler(Fail)

	var specReporters []Reporter
	junitReporter := reporters.NewJUnitReporter("/testresults/junit_basic.xml")
	specReporters = append(specReporters, junitReporter)
	RunSpecsWithDefaultAndCustomReporters(t, "Torpedo : Backup", specReporters)
}

var _ = BeforeSuite(func() {
	InitInstance()
})

// This test performs basic test of starting an application and destroying it (along with storage)
var _ = Describe("{BackupSetup}", func() {
	//var contexts []*scheduler.Context

	It("has to validate that backup completes even after killing storage driver", func() {
		/*Step("Deploy applications", func() {
			contexts = make([]*scheduler.Context, 0)

			for i := 0; i < Inst().ScaleFactor; i++ {
				contexts = append(contexts, ScheduleApplications(fmt.Sprintf("simplebackup-%d", i))...)

			}
			ValidateApplications(contexts)
		})
		*/

		CreateOrganization(orgID)

		CreateCloudCredential(CredName, orgID)

		CreateBackupLocation(BLocationName, orgID, CredName)

		CreateCluster(ClusterName, orgID, CredName)
	})
})

var _ = AfterSuite(func() {
	//PerformSystemCheck()
	//ValidateCleanup()
})

func CreateOrganization(name string) {

	Step(fmt.Sprintf("Create organization [%s]", name), func() {
		backupDriver := Inst().Backup
		Expect(backupDriver).NotTo(BeNil())
		//logrus.Infof("backup driver name :%v", Inst().Backup.String())
		metadata := &api.CreateMetadata{
			Name: name,
		}
		createOrgRequest := &api.OrganizationCreateRequest{
			CreateMetadata: metadata,
		}
		_, err := backupDriver.CreateOrganization(createOrgRequest)
		Expect(err).NotTo(HaveOccurred())
		// TODO: validate createOrgResponse also
	})
}

func CreateCloudCredential(name string, orgID string) {

	Step(fmt.Sprintf("Create cloud credential [%s] in org [%s]", name, orgID), func() {
		backupDriver := Inst().Backup

		// TODO: add separate function to return cred object based on type
		id := os.Getenv("AWS_ACCESS_KEY_ID")
		Expect(id).NotTo(Equal(""))

		secret := os.Getenv("AWS_SECRET_ACCESS_KEY")
		Expect(secret).NotTo(Equal(""))

		awsConfig := &api.AWSConfig{
			AccessKey: id,
			SecretKey: secret,
		}
		credConfig := &api.CloudCredentialInfo_AwsConfig{
			AwsConfig: awsConfig,
		}
		metadata := &api.CreateMetadata{
			Name:  name,
			OrgId: orgID,
		}
		credInfo := &api.CloudCredentialInfo{
			Type:   api.CloudCredentialInfo_AWS,
			Config: credConfig,
		}
		credCreateRequest := &api.CloudCredentialCreateRequest{
			CreateMetadata:  metadata,
			CloudCredential: credInfo,
		}
		_, err := backupDriver.CreateCloudCredential(credCreateRequest)
		Expect(err).NotTo(HaveOccurred())
		// TODO: validate CreateCloudCredentialResponse also
	})

}

func CreateBackupLocation(name string, orgID string, cloudCred string) {

	Step(fmt.Sprintf("Create backup location [%s] in org [%s]", name, orgID), func() {
		backupDriver := Inst().Backup

		// TODO:
		path := "abc"
		encryptionKey := "abc"

		metadata := &api.CreateMetadata{
			Name:  name,
			OrgId: orgID,
		}
		bLocationInfo := &api.BackupLocationInfo{
			Path:            path,
			EncryptionKey:   encryptionKey,
			CloudCredential: cloudCred,
		}
		bLocationCreateReq := &api.BackupLocationCreateRequest{
			CreateMetadata: metadata,
			BackupLocation: bLocationInfo,
		}

		_, err := backupDriver.CreateBackupLocation(bLocationCreateReq)
		Expect(err).NotTo(HaveOccurred())
		// TODO: validate createBackupLocationResponse also
	})
}

func CreateCluster(name string, orgID string, cloudCred string) {

	Step(fmt.Sprintf("Create cluster [%s] in org [%s]", name, orgID), func() {
		backupDriver := Inst().Backup

		// TODO: Add kubeconfig as k8s secret mount
		kubeconfig := "abc"

		metadata := &api.CreateMetadata{
			Name:  name,
			OrgId: orgID,
		}

		clusterInfo := &api.ClusterInfo{
			Kubeconfig:      kubeconfig,
			CloudCredential: cloudCred,
		}
		clusterCreateReq := &api.ClusterCreateRequest{
			CreateMetadata: metadata,
			Cluster:        clusterInfo,
		}
		_, err := backupDriver.CreateCluster(clusterCreateReq)
		Expect(err).NotTo(HaveOccurred())
		// TODO: validate createClusterResponse also
	})
}

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	ParseFlags()
	os.Exit(m.Run())
}
