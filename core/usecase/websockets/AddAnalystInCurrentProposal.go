package websockets

import (
	"github.com/gtkmk/finder_api/core/domain/proposalAnalysis"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type AddAnalystInCurrentProposal struct {
	proposalAnalysisDatabase repositories.ProposalAnalysisRepositoryInterface
}

func NewAddAnalystInCurrentProposal(
	proposalDatabase repositories.ProposalAnalysisRepositoryInterface,
) *AddAnalystInCurrentProposal {
	return &AddAnalystInCurrentProposal{proposalDatabase}
}

func (addAnalystInCurrentProposal *AddAnalystInCurrentProposal) Execute(analysis *proposalAnalysisDomain.ProposalAnalysis) ([]map[string]interface{}, error) {
	if err := addAnalystInCurrentProposal.proposalAnalysisDatabase.AddAnalyst(analysis); err != nil {
		return nil, err
	}

	return addAnalystInCurrentProposal.proposalAnalysisDatabase.FindAllAnalystsInCurrentProposal(analysis.ProposalId)
}
