package workers

import (
	"github.com/gtkmk/finder_api/core/domain/permissionDomain"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type UserReportWorkerTask struct {
	database   repositories.UserRepository
	unityIds   []any
	userId     string
	groupLayer float64
	offset     int
	limit      int
	reportData *[]map[string]interface{}
}

func NewUserReportWorkerTask(
	database repositories.UserRepository,
	unityIdsList []any,
	userId string,
	groupLayer float64,
	offset int,
	limit int,
	reportData *[]map[string]interface{},
) UserReportWorkerTask {
	return UserReportWorkerTask{
		database:   database,
		unityIds:   unityIdsList,
		userId:     userId,
		groupLayer: groupLayer,
		offset:     offset,
		limit:      limit,
		reportData: reportData,
	}
}

func (userReportWorkerTask UserReportWorkerTask) ExecuteWorker() error {
	var err error
	var data []map[string]interface{}

	if userReportWorkerTask.groupLayer >= permissionDomain.ManagersLayerLimitConst {
		data, err = userReportWorkerTask.database.FindUsersForReport(
			userReportWorkerTask.unityIds,
			"",
			userReportWorkerTask.userId,
			userReportWorkerTask.offset,
			userReportWorkerTask.limit,
		)
	} else {
		data, err = userReportWorkerTask.database.FindUsersForReport(
			userReportWorkerTask.unityIds,
			userReportWorkerTask.userId,
			userReportWorkerTask.userId,
			userReportWorkerTask.offset,
			userReportWorkerTask.limit,
		)
	}

	*userReportWorkerTask.reportData = append(*userReportWorkerTask.reportData, data...)

	return err
}
