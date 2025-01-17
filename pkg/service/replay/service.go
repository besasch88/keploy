package replay

import (
	"context"
	"time"

	"go.keploy.io/server/v2/pkg/models"
)

type Instrumentation interface {
	//Setup prepares the environment for the recording
	Setup(ctx context.Context, cmd string, opts models.SetupOptions) (uint64, error)
	//Hook will load hooks and start the proxy server.
	Hook(ctx context.Context, id uint64, opts models.HookOptions) error
	MockOutgoing(ctx context.Context, id uint64, opts models.OutgoingOptions) error
	// SetMocks Allows for setting mocks between test runs for better filtering and matching
	SetMocks(ctx context.Context, id uint64, filtered []*models.Mock, unFiltered []*models.Mock) error
	// GetConsumedMocks to log the names of the mocks that were consumed during the test run of failed test cases
	GetConsumedMocks(ctx context.Context, id uint64) ([]string, error)
	// Run is blocking call and will execute until error
	Run(ctx context.Context, id uint64, opts models.RunOptions) models.AppError

	GetContainerIP(ctx context.Context, id uint64) (string, error)
}

type Service interface {
	Start(ctx context.Context) error
	Instrument(ctx context.Context) (*InstrumentState, error)
	GetNextTestRunID(ctx context.Context) (string, error)
	GetAllTestSetIDs(ctx context.Context) ([]string, error)
	RunTestSet(ctx context.Context, testSetID string, testRunID string, appID uint64, serveTest bool) (models.TestSetStatus, error)
	GetTestSetStatus(ctx context.Context, testRunID string, testSetID string) (models.TestSetStatus, error)
	RunApplication(ctx context.Context, appID uint64, opts models.RunOptions) models.AppError
	Normalize(ctx context.Context) error
	DenoiseTestCases(ctx context.Context, testSetID string, noiseParams []*models.NoiseParams) ([]*models.NoiseParams, error)
	NormalizeTestCases(ctx context.Context, testRun string, testSetID string, selectedTestCaseIDs []string, testResult []models.TestResult) error
	DeleteTests(ctx context.Context, testSetID string, testCaseIDs []string) error
	DeleteTestSet(ctx context.Context, testSetID string) error
}

type TestDB interface {
	GetAllTestSetIDs(ctx context.Context) ([]string, error)
	GetTestCases(ctx context.Context, testSetID string) ([]*models.TestCase, error)
	UpdateTestCase(ctx context.Context, testCase *models.TestCase, testSetID string) error
	DeleteTests(ctx context.Context, testSetID string, testCaseIDs []string) error
	DeleteTestSet(ctx context.Context, testSetID string) error
}

type MockDB interface {
	GetFilteredMocks(ctx context.Context, testSetID string, afterTime time.Time, beforeTime time.Time) ([]*models.Mock, error)
	GetUnFilteredMocks(ctx context.Context, testSetID string, afterTime time.Time, beforeTime time.Time) ([]*models.Mock, error)
	UpdateMocks(ctx context.Context, testSetID string, mockNames map[string]bool) error
}

type ReportDB interface {
	GetAllTestRunIDs(ctx context.Context) ([]string, error)
	GetTestCaseResults(ctx context.Context, testRunID string, testSetID string) ([]models.TestResult, error)
	GetReport(ctx context.Context, testRunID string, testSetID string) (*models.TestReport, error)
	InsertTestCaseResult(ctx context.Context, testRunID string, testSetID string, result *models.TestResult) error
	InsertReport(ctx context.Context, testRunID string, testSetID string, testReport *models.TestReport) error
}

type Config interface {
	Read(ctx context.Context, testSetID string) (*models.TestSet, error)
	Write(ctx context.Context, testSetID string, testSet *models.TestSet) error
}

type Telemetry interface {
	TestSetRun(success int, failure int, testSet string, runStatus string)
	TestRun(success int, failure int, testSets int, runStatus string)
	MockTestRun(utilizedMocks int)
}

// RequestMockHandler defines an interface for implementing hooks that extend and customize
// the behavior of request simulations and test workflows. This interface allows for
// detailed control over various stages of the testing process, including request simulation,
// test status processing, and post-test actions.
type RequestMockHandler interface {
	SimulateRequest(ctx context.Context, appID uint64, tc *models.TestCase, testSetID string) (*models.HTTPResp, error)
	ProcessTestRunStatus(ctx context.Context, status bool, testSetID string)
	FetchMockName() string
	ProcessMockFile(ctx context.Context, testSetID string)
	AfterTestHook(ctx context.Context, testRunID, testSetID string, totalTestSets int) (*models.TestReport, error)
}

type InstrumentState struct {
	AppID      uint64
	HookCancel context.CancelFunc
}

type MockAction string

// MockAction constants define the possible actions that can be taken on a mocking.
const (
	Start  MockAction = "start"
	Update MockAction = "update"
)
