package workers

import (
	"github.com/gtkmk/finder_api/core/domain/proposalInstallments"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type InstallmentWorkerTask struct {
	installment proposalInstallmentsDomain.ProposalInstallments
	database    repositories.ProposalRepositoryInterface
}

func NewInstallmentWorkerTask(
	installment proposalInstallmentsDomain.ProposalInstallments,
	database repositories.ProposalRepositoryInterface,
) InstallmentWorkerTask {
	return InstallmentWorkerTask{
		installment,
		database,
	}
}

func (installmentWorkerTask InstallmentWorkerTask) ExecuteWorker() error {
	return installmentWorkerTask.database.CreateProposalInstallment(&installmentWorkerTask.installment)
}
