package export

import (
	"context"
	"encoding/csv"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"io/ioutil"
	"math"
	"os"
	"sync"
)

// ExportCsv fetches data by page and saves it in csv file.
func ExportCsv(
	exportFields []string,
	separator rune,
	dataFetcher DataFetcher,
) (resFileName string, resErr error) {
	if len(exportFields) == 0 {
		return "", fmt.Errorf("exportFields is empty")
	}

	file, err := ioutil.TempFile("", "export.*.csv")
	if err != nil {
		return "", err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			return
		}

		if resErr != nil {
			err := os.Remove(file.Name())
			if err != nil {
				return
			}
		}
	}()

	w := csv.NewWriter(file)
	if separator != 0 {
		w.Comma = separator
	}

	err = w.WriteAll([][]string{exportFields})
	if err != nil {
		return "", err
	}

	var limit int64 = 1000
	data, totalCount, err := dataFetcher(1, limit)
	if err != nil {
		return "", err
	}

	if totalCount > 0 {
		err = utils.ExportCsv(w, data, exportFields)
		if err != nil {
			return "", err
		}

		pageCount := int64(math.Ceil(float64(totalCount) / float64(limit)))
		if pageCount > 1 {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			pageCh := make(chan int64)
			dataCh := runExportWorkers(ctx, pageCh, dataFetcher, limit)

			go func() {
				defer close(pageCh)
				for p := int64(2); p <= pageCount; p++ {
					pageCh <- p
				}
			}()

			var err error
			dataByPage := make(map[int64]interface{})
			page := int64(2)
			for res := range dataCh {
				if res.Err != nil {
					err = res.Err
					continue
				}

				if res.Page == page {
					err = utils.ExportCsv(w, res.Data, exportFields)
					if err != nil {
						return "", err
					}

					for p := page + 1; p <= pageCount; p++ {
						if d, ok := dataByPage[p]; ok {
							err = utils.ExportCsv(w, d, exportFields)
							if err != nil {
								return "", err
							}

							delete(dataByPage, p)
						} else {
							page = p
							break
						}
					}
				} else {
					dataByPage[res.Page] = res.Data
				}
			}

			if err != nil {
				return "", err
			}
		}
	}

	return file.Name(), nil
}

type workerResult struct {
	Page int64
	Data interface{}
	Err  error
}

// runExportWorkers starts workers. Each worker fetches specific page and sends csv data to output channel.
func runExportWorkers(
	ctx context.Context,
	in <-chan int64,
	dataFetcher DataFetcher,
	limit int64,
) <-chan workerResult {
	out := make(chan workerResult)

	go func() {
		defer close(out)

		wg := sync.WaitGroup{}

		for i := 0; i < 10; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()

				for {
					select {
					case <-ctx.Done():
						return
					case page, ok := <-in:
						if !ok {
							return
						}

						data, _, err := dataFetcher(page, limit)
						if err != nil {
							out <- workerResult{
								Page: page,
								Err:  err,
							}
							return
						}

						out <- workerResult{
							Page: page,
							Data: data,
						}
					}
				}
			}()
		}

		wg.Wait()
	}()

	return out
}
