package internal

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Repository interface {
	GetBeforeEod(ctx context.Context) (list BeforeEodList, err error)
	UpdateAfterEod(ctx context.Context, data AfterEodList) error
}

type repository struct {
	pathFolder string
}

func NewRepository(pathFolder string) Repository {
	return &repository{
		pathFolder: pathFolder,
	}
}

func (r *repository) GetBeforeEod(ctx context.Context) (list BeforeEodList, err error) {
	filePath := filepath.Join(r.pathFolder, string(BeforeEodFile))
	f, err := os.Open(filePath)
	if err != nil {
		log.Println("[Repository][GetBeforeEod] Unable to read input file "+filePath, err)
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Println("[Repository][GetBeforeEod] Unable to parse file as CSV for "+filePath, err)
		return nil, err
	}

	if len(records) < 7 {
		err = errors.New("data is not valid")
		log.Println("[Repository][GetBeforeEod] ", err)
		return nil, err
	}

	for i := 1; i < len(records); i++ {
		record := records[i]
		beforeEod := BeforeEod{}
		for j := 0; j < len(record); j++ {
			switch j {
			case 0:
				beforeEod.ID, _ = strconv.ParseInt(record[j], 10, 64)
			case 1:
				beforeEod.Name = record[j]
			case 2:
				beforeEod.Age, _ = strconv.ParseInt(record[j], 10, 64)
			case 3:
				beforeEod.Balanced, _ = strconv.ParseFloat(record[j], 64)
			case 4:
				beforeEod.PreviousBalanced, _ = strconv.ParseFloat(record[j], 64)
			case 5:
				beforeEod.AverageBalanced, _ = strconv.ParseFloat(record[j], 64)
			case 6:
				beforeEod.FreeTransfer, _ = strconv.ParseInt(record[j], 10, 64)
			}
		}
		list = append(list, &beforeEod)
	}
	return
}

func (r *repository) UpdateAfterEod(ctx context.Context, data AfterEodList) error {
	filePath := filepath.Join(r.pathFolder, string(AfterEodFile))
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("[Repository][UpdateAfterEod]failed to open file", err)
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.Comma = ';'
	defer w.Flush()
	var records [][]string
	header := []string{
		string(ID), string(Name), string(Age), string(Balanced),
		string(ProcessTwoB), string(ProcessThree), string(PreviousBalanced),
		string(AverageBalanced), string(ProcessOne), string(FreeTransfer), string(ProcessTwoA),
	}
	records = append(records, header)
	for _, d := range data {
		row := []string{}
		row = append(row, fmt.Sprint(d.ID))
		row = append(row, d.Name)
		row = append(row, fmt.Sprint(d.Age))
		row = append(row, fmt.Sprint(d.Balanced))
		processTwoB := ""
		if d.ProcessTwoB != nil {
			processTwoB = fmt.Sprint(*d.ProcessTwoB)
		}
		row = append(row, processTwoB)
		processThree := ""
		if d.ProcessThree != nil {
			processThree = fmt.Sprint(*d.ProcessThree)
		}
		row = append(row, processThree)
		row = append(row, fmt.Sprint(d.PreviousBalanced))
		row = append(row, fmt.Sprint(d.AverageBalanced))
		row = append(row, fmt.Sprint(d.ProcessOne))
		row = append(row, fmt.Sprint(d.FreeTransfer))
		processTwoA := ""
		if d.ProcessTwoA != nil {
			processTwoA = fmt.Sprint(*d.ProcessTwoA)
		}
		row = append(row, processTwoA)
		records = append(records, row)
	}
	err = w.WriteAll(records)
	if err != nil {
		log.Println("[Repository][UpdateAfterEod]failed to write file", err)
	}
	return err
}
