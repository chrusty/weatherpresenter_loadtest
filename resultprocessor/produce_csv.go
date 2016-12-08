package resultprocessor

import (
	fmt "fmt"
	strconv "strconv"
)

func ProduceCSV() {
	// Get a lock on the report (because it could be changed in the background if we call this when the tests are still running):
	reportMutex.Lock()
	defer reportMutex.Unlock()

	// Print a CSV header:
	fmt.Println(`"Request Time","Test Name","Machine Address","Concurrency","Duration (ns)","HTTP Response Code","Timed Out?","Error?"`)

	// Add records:
	for _, result := range report.Results {
		fmt.Printf("\"%s\",\"%s\",\"%s\",%d,%d,%d,%s,\"%s\"\n", result.RequestTime, result.TestName, result.MachineAddress, myConfig.Concurrency, result.Duration, result.ResponseCode, strconv.FormatBool(result.TimedOut), result.Error)
	}
}
