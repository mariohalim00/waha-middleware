package job

import "waha-job-processing/internal/service/waha"

func PreprocessJobs() {

	waha.GetAllActiveSessions()
}
