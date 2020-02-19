package portworx

import (
	"context"
	"fmt"

	api "github.com/portworx/px-backup-api/pkg/apis/v1"
	"github.com/portworx/sched-ops/k8s/core"
	"github.com/portworx/torpedo/drivers/backup"
	"github.com/portworx/torpedo/drivers/node"
	"github.com/portworx/torpedo/drivers/scheduler"
	"github.com/portworx/torpedo/drivers/volume"
	"github.com/portworx/torpedo/drivers/volume/portworx/schedops"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	driverName            = "pxb"
	pxbRestPort           = 10001
	defaultPxbServicePort = 10002
	pxbServiceName        = "px-backup"
	pxbNameSpace          = "px-backup"
	schedulerDriverName   = "k8s"
	nodeDriverName        = "ssh"
	volumeDriverName      = "pxd"
)

type portworx struct {
	clusterManager         api.ClusterClient
	backupLocationManager  api.BackupLocationClient
	cloudCredentialManager api.CloudCredentialClient
	backupManger           api.BackupClient
	restoreManager         api.RestoreClient
	backupScheduleManager  api.BackupScheduleClient
	schedulePolicyManager  api.SchedulePolicyClient
	organizationManager    api.OrganizationClient
	healthManager          api.HealthClient

	schedulerDriver scheduler.Driver
	nodeDriver      node.Driver
	volumeDriver    volume.Driver
	schedOps        schedops.Driver
	refreshEndpoint bool
	token           string
}

func (p *portworx) String() string {
	return driverName
}

func (p *portworx) Init(schedulerDriverName string, nodeDriverName string, volumeDriverName string, token string) error {
	var err error

	logrus.Infof("using portworx px-backup driver under scheduler: %v", schedulerDriverName)

	p.nodeDriver, err = node.Get(nodeDriverName)
	if err != nil {
		return err
	}
	p.token = token

	p.schedulerDriver, err = scheduler.Get(schedulerDriverName)
	if err != nil {
		return fmt.Errorf("Error getting scheduler driver %v: %v", schedulerDriverName, err)
	}

	p.volumeDriver, err = volume.Get(volumeDriverName)
	if err != nil {
		return fmt.Errorf("Error getting volume driver %v: %v", volumeDriverName, err)
	}

	if err = p.setDriver(pxbServiceName, pxbNameSpace); err != nil {
		return fmt.Errorf("Error setting px-backup endpoint: %v", err)
	}

	return err

}

func (p *portworx) constructURL(ip string) string {
	return fmt.Sprintf("%s:%d", ip, defaultPxbServicePort)
}

func (p *portworx) testAndSetEndpoint(endpoint string) error {
	pxEndpoint := p.constructURL(endpoint)
	conn, err := grpc.Dial(pxEndpoint, grpc.WithInsecure())
	if err != nil {
		logrus.Errorf("unable to get grpc connection: %v", err)
		return err
	}

	p.healthManager = api.NewHealthClient(conn)
	_, err = p.healthManager.Status(context.Background(), &api.HealthStatusRequest{})
	if err != nil {
		logrus.Errorf("HealthManager API error: %v", err)
		return err
	}

	p.clusterManager = api.NewClusterClient(conn)
	p.backupLocationManager = api.NewBackupLocationClient(conn)
	p.cloudCredentialManager = api.NewCloudCredentialClient(conn)
	p.backupManger = api.NewBackupClient(conn)
	p.restoreManager = api.NewRestoreClient(conn)
	p.backupScheduleManager = api.NewBackupScheduleClient(conn)
	p.schedulePolicyManager = api.NewSchedulePolicyClient(conn)
	p.organizationManager = api.NewOrganizationClient(conn)

	logrus.Infof("Using %v as endpoint for portworx backup driver", pxEndpoint)

	return err
}

func (p *portworx) GetServiceEndpoint(serviceName string, nameSpace string) (string, error) {
	svc, err := core.Instance().GetService(serviceName, nameSpace)
	if err == nil {
		return svc.Spec.ClusterIP, nil
	}
	return "", err
}

func (p *portworx) setDriver(serviceName string, nameSpace string) error {
	var err error
	var endpoint string

	endpoint, err = p.GetServiceEndpoint(serviceName, nameSpace)
	if err == nil && endpoint != "" {
		if err = p.testAndSetEndpoint(endpoint); err == nil {
			p.refreshEndpoint = false
			return nil
		}
	} else if err != nil && len(node.GetWorkerNodes()) == 0 {
		return err
	}

	p.refreshEndpoint = true
	for _, n := range node.GetWorkerNodes() {
		for _, addr := range n.Addresses {
			if err = p.testAndSetEndpoint(addr); err == nil {
				return nil
			}
		}
	}

	return fmt.Errorf("failed to get endpoint for portworx backup driver: %v", err)
}

func (p *portworx) getOrganizationManager() api.OrganizationClient {
	if p.refreshEndpoint {
		p.setDriver(pxbServiceName, pxbNameSpace)
	}
	return p.organizationManager
}

func (p *portworx) getClusterManager() api.ClusterClient {
	if p.refreshEndpoint {
		p.setDriver(pxbServiceName, pxbNameSpace)
	}
	return p.clusterManager
}

func (p *portworx) getBackupLocationManager() api.BackupLocationClient {
	if p.refreshEndpoint {
		p.setDriver(pxbServiceName, pxbNameSpace)
	}
	return p.backupLocationManager
}

func (p *portworx) getCloudCredentialManager() api.CloudCredentialClient {
	if p.refreshEndpoint {
		p.setDriver(pxbServiceName, pxbNameSpace)
	}
	return p.cloudCredentialManager
}

func (p *portworx) getBackupManager() api.BackupClient {
	if p.refreshEndpoint {
		p.setDriver(pxbServiceName, pxbNameSpace)
	}
	return p.backupManger
}

func (p *portworx) getRestoreManager() api.RestoreClient {
	if p.refreshEndpoint {
		p.setDriver(pxbServiceName, pxbNameSpace)
	}
	return p.restoreManager
}

func (p *portworx) getBackupScheduleManager() api.BackupScheduleClient {
	if p.refreshEndpoint {
		p.setDriver(pxbServiceName, pxbNameSpace)
	}
	return p.backupScheduleManager
}

func (p *portworx) getSchedulePolicyManager() api.SchedulePolicyClient {
	if p.refreshEndpoint {
		p.setDriver(pxbServiceName, pxbNameSpace)
	}
	return p.schedulePolicyManager
}

func (p *portworx) CreateOrganization(req *api.OrganizationCreateRequest) (*api.OrganizationCreateResponse, error) {
	org := p.getOrganizationManager()
	if req != nil {
		return org.Create(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) EnumerateOrganization() (*api.OrganizationEnumerateResponse, error) {
	org := p.getOrganizationManager()
	req := &api.OrganizationEnumerateRequest{}

	return org.Enumerate(context.Background(), req)
}

func (p *portworx) CreateCloudCredential(req *api.CloudCredentialCreateRequest) (*api.CloudCredentialCreateResponse, error) {
	cc := p.getCloudCredentialManager()
	if req != nil {
		return cc.Create(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) InspectCloudCredential(req *api.CloudCredentialInspectRequest) (*api.CloudCredentialInspectResponse, error) {
	cc := p.getCloudCredentialManager()
	if req != nil {
		return cc.Inspect(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) EnumerateCloudCredential(req *api.CloudCredentialEnumerateRequest) (*api.CloudCredentialEnumerateResponse, error) {
	cc := p.getCloudCredentialManager()
	if req != nil {
		cc.Enumerate(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) DeleteCloudCredential(req *api.CloudCredentialDeleteRequest) (*api.CloudCredentialDeleteResponse, error) {
	cc := p.getCloudCredentialManager()
	if req != nil {
		return cc.Delete(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) CreateCluster(req *api.ClusterCreateRequest) (*api.ClusterCreateResponse, error) {
	cluster := p.getClusterManager()
	if req != nil {
		return cluster.Create(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) InspectCluster(req *api.ClusterInspectRequest) (*api.ClusterInspectResponse, error) {
	cluster := p.getClusterManager()
	if req != nil {
		return cluster.Inspect(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) EnumerateCluster(req *api.ClusterEnumerateRequest) (*api.ClusterEnumerateResponse, error) {
	cluster := p.getClusterManager()
	if req != nil {
		return cluster.Enumerate(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) DeleteCluster(req *api.ClusterDeleteRequest) (*api.ClusterDeleteResponse, error) {
	cluster := p.getClusterManager()
	if req != nil {
		return cluster.Delete(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) CreateBackupLocation(req *api.BackupLocationCreateRequest) (*api.BackupLocationCreateResponse, error) {
	bl := p.getBackupLocationManager()
	if req != nil {
		return bl.Create(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) EnumerateBackupLocation(req *api.BackupLocationEnumerateRequest) (*api.BackupLocationEnumerateResponse, error) {
	bl := p.getBackupLocationManager()
	if req != nil {
		return bl.Enumerate(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) InspectBackupLocation(req *api.BackupLocationInspectRequest) (*api.BackupLocationInspectResponse, error) {
	bl := p.getBackupLocationManager()
	if req != nil {
		return bl.Inspect(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) DeleteBackupLocation(req *api.BackupLocationDeleteRequest) (*api.BackupLocationDeleteResponse, error) {
	bl := p.getBackupLocationManager()
	if req != nil {
		return bl.Delete(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) CreateBackup(req *api.BackupCreateRequest) (*api.BackupCreateResponse, error) {
	backup := p.getBackupManager()
	if req != nil {
		return backup.Create(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) EnumerateBackup(req *api.BackupEnumerateRequest) (*api.BackupEnumerateResponse, error) {
	backup := p.getBackupManager()
	if req != nil {
		return backup.Enumerate(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) InspectBackup(req *api.BackupInspectRequest) (*api.BackupInspectResponse, error) {
	backup := p.getBackupManager()
	if req != nil {
		return backup.Inspect(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func (p *portworx) DeleteBackup(req *api.BackupDeleteRequest) (*api.BackupDeleteResponse, error) {
	backup := p.getBackupManager()
	if req != nil {
		return backup.Delete(context.Background(), req)
	}
	return nil, fmt.Errorf("Request is nil")
}

func init() {
	backup.Register(driverName, &portworx{})
}
