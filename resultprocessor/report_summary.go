package resultprocessor

import (
	logrus "github.com/Sirupsen/logrus"
)

func ReportSummary() {
	// Get a lock on the report (because it could be changed in the background if we call this when the tests are still running):
	reportMutex.Lock()
	defer reportMutex.Unlock()

	// Log the report:
	log.WithFields(logrus.Fields{"tests_run": len(report.Results)}).Info("Report")
	log.WithFields(logrus.Fields{"test_errors": report.Errors}).Info("Report")
	log.WithFields(logrus.Fields{"test_timeouts": report.TimeOuts}).Info("Report")
	log.WithFields(logrus.Fields{"test_successes": report.Successes}).Info("Report")
	log.WithFields(logrus.Fields{"duration_minimum": report.MinimumDuration}).Info("Report")
	log.WithFields(logrus.Fields{"duration_average": report.AverageDuration}).Info("Report")
	log.WithFields(logrus.Fields{"duration_maximum": report.MaximumDuration}).Info("Report")
}
