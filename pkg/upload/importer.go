package upload

import (
	"encoding/csv"
)

type Data struct {
	Names []string
	Data  [][]string
}

func (d *Data) Size() int {
	return len(d.Names)
}

type Importer interface {
	Import(*Upload) (*Data, error)
}

type CSVImporter struct{}

func (i *CSVImporter) Import(up *Upload) (*Data, error) {
	csvReader := csv.NewReader(up.ImportReader)
	csvReader.TrimLeadingSpace = true

	allRows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	data := &Data{
		Names: allRows[0],
		Data:  allRows[1:],
	}
	return data, nil
}
