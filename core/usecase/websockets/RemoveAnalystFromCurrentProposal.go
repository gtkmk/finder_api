package websockets

import (
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type RemoveAnalystFromCurrentProposal struct {
	proposalAnalysisDatabase repositories.ProposalAnalysisRepositoryInterface
}

func NewRemoveAnalystFromCurrentProposal(
	proposalDatabase repositories.ProposalAnalysisRepositoryInterface,
) *RemoveAnalystFromCurrentProposal {
	return &RemoveAnalystFromCurrentProposal{proposalDatabase}
}

func (removeAnalystFromCurrentProposal *RemoveAnalystFromCurrentProposal) Execute(analystId string) error {
	return removeAnalystFromCurrentProposal.proposalAnalysisDatabase.DeleteAnalyst(analystId)
}
