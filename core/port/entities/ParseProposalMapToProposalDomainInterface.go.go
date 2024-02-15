package proposalEntityInterface

import "github.com/gtkmk/finder_api/core/domain/proposalDomain"

type ProposalEntityInterface interface {
	MountProposalData(proposalMap []map[string]interface{}) (map[string]interface{}, error)
	VerifyIfProposalCanBeApprovedOrFormalized(proposal map[string]interface{}, proposalApprovedAndFormalizedInfos *proposalDomain.ProposalApprovedAndFormalized) error
}
