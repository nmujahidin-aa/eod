package internal

type ColumnName string

var (
	ID               ColumnName = "id"
	Name             ColumnName = "Nama"
	Age              ColumnName = "Age"
	Balanced         ColumnName = "Balanced"
	PreviousBalanced ColumnName = "Previous Balanced"
	AverageBalanced  ColumnName = "Average Balanced"
	FreeTransfer     ColumnName = "Free Transfer"
	ProcessOne       ColumnName = "No 1 Thread-No"
	ProcessTwoA      ColumnName = "No 2a Thread-No"
	ProcessTwoB      ColumnName = "No 2b Thread-No"
	ProcessThree     ColumnName = "No 3 Thread-No"
)

type AfterEod struct {
	ID               int64
	Name             string
	Age              int64
	Balanced         float64
	ProcessTwoB      *int64
	ProcessThree     *int64
	PreviousBalanced float64
	AverageBalanced  float64
	ProcessOne       int64
	FreeTransfer     int64
	ProcessTwoA      *int64
}
type AfterEodList []*AfterEod
type BeforeEod struct {
	ID               int64
	Name             string
	Age              int64
	Balanced         float64
	PreviousBalanced float64
	AverageBalanced  float64
	FreeTransfer     int64
}
type BeforeEodList []*BeforeEod

func NewAfterEod(beforeEod BeforeEod) *AfterEod {
	return &AfterEod{
		ID:               beforeEod.ID,
		Name:             beforeEod.Name,
		Age:              beforeEod.Age,
		Balanced:         beforeEod.Balanced,
		PreviousBalanced: beforeEod.PreviousBalanced,
		AverageBalanced:  beforeEod.AverageBalanced,
		FreeTransfer:     beforeEod.FreeTransfer,
	}
}

type DataName string

var (
	BeforeEodFile DataName = "Before Eod.csv"
	AfterEodFile  DataName = "After Eod.csv"
)
