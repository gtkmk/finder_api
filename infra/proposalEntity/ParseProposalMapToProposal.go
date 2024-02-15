package proposalEntity

import (
	"github.com/gtkmk/finder_api/core/domain/helper"
	"reflect"
	"strconv"

	"github.com/gtkmk/finder_api/core/domain/contractDomain"
	"github.com/gtkmk/finder_api/core/domain/personDomain"
	"github.com/gtkmk/finder_api/core/domain/proposalDomain"
	
	proposalEntityInterface "github.com/gtkmk/finder_api/core/port/entities"
)

const BaseProposalPrefixConst = "base_proposal_"

type ProposalEntity struct{}

func NewProposalEntity() proposalEntityInterface.ProposalEntityInterface {
	return &ProposalEntity{}
}

func (proposalEntity *ProposalEntity) MountProposalData(proposalMap []map[string]interface{}) (map[string]interface{}, error) {
	var (
		guarantorDocuments       []map[string]interface{}
		clientDocuments          []map[string]interface{}
		clientCpf                string
		clientId                 string
		clientName               string
		clientEmail              string
		clientCellphoneNumber    string
		clientTelephoneNumber    string
		clientBirthDate          interface{}
		guarantorCpf             string
		guarantorId              string
		guarantorName            string
		guarantorEmail           string
		guarantorCellphoneNumber string
		enterpriseValue          float64
		maximumInstallmentValue  float64
		guarantorTelephoneNumber interface{}
		guarantorCompleteAddress interface{}
		guarantorComplement      interface{}
		guarantorNeighborhood    interface{}
		guarantorZipCode         interface{}
		guarantorCity            interface{}
		guarantorBirthDate       interface{}
		guarantorUf              interface{}
		guarantorAddressCode     interface{}
		proposalApproved         bool
		proposalFormalized       bool
		databasePrefix           string
	)

	if proposalMap[0]["operation"] == proposalDomain.ProposalOperationDifferenceConst {
		databasePrefix = BaseProposalPrefixConst
	}

	for _, rangeProposal := range proposalMap {
		if rangeProposal["person_type"] == personDomain.ClientTypeConst {
			var documentType string
			var documentIsPending bool

			if rangeProposal["document_type"] != nil {
				documentType = rangeProposal["document_type"].(string)
			}

			if rangeProposal["document_is_pending"] != nil {
				documentIsPending = rangeProposal["document_is_pending"].(int64) > 0
			}

			if documentType != "" {
				clientDocuments = append(clientDocuments, map[string]interface{}{
					"type":       documentType,
					"link":       rangeProposal["document_path"],
					"is_pending": documentIsPending,
					"id":         rangeProposal["document_id"],
				})
			}

			clientCpf = rangeProposal["cpf"].(string)
			clientName = rangeProposal["name"].(string)
			clientId = rangeProposal["person_id"].(string)
			clientEmail = rangeProposal["email"].(string)
			clientCellphoneNumber = rangeProposal["cellphone_number"].(string)
			clientTelephoneNumber = rangeProposal["telephone_number"].(string)
			clientBirthDate = rangeProposal["birth_date"]
		} else if rangeProposal["person_type"] == personDomain.GuarantorTypeConst {
			var documentType string
			var documentIsPending bool

			if rangeProposal["document_type"] != nil {
				documentType = rangeProposal["document_type"].(string)
			}

			if rangeProposal["document_is_pending"] != nil {
				documentIsPending = rangeProposal["document_is_pending"].(int64) > 0
			}

			if documentType != "" {
				guarantorDocuments = append(guarantorDocuments, map[string]interface{}{
					"type":       documentType,
					"link":       rangeProposal["document_path"],
					"is_pending": documentIsPending,
					"id":         rangeProposal["document_id"],
				})
			}

			guarantorCpf = rangeProposal["cpf"].(string)
			guarantorId = rangeProposal["person_id"].(string)
			guarantorName = rangeProposal["name"].(string)
			guarantorEmail = rangeProposal["email"].(string)
			guarantorCellphoneNumber = rangeProposal["cellphone_number"].(string)
			guarantorCompleteAddress = rangeProposal["complete_address"]
			guarantorComplement = rangeProposal["complement"]
			guarantorNeighborhood = rangeProposal["neighborhood"]
			guarantorZipCode = rangeProposal["zip_code"]
			guarantorCity = rangeProposal["city"]
			guarantorBirthDate = rangeProposal["birth_date"]
			guarantorUf = rangeProposal["uf"]
			guarantorTelephoneNumber = rangeProposal["telephone_number"]
			guarantorAddressCode = rangeProposal["address_code"]
		}
	}

	guarantorInfo := map[string]interface{}{
		"cpf":                        guarantorCpf,
		"name":                       guarantorName,
		"email":                      guarantorEmail,
		"cellphone_number":           guarantorCellphoneNumber,
		"documents":                  guarantorDocuments,
		"id":                         guarantorId,
		"guarantor_complete_address": guarantorCompleteAddress,
		"guarantor_complement":       guarantorComplement,
		"guarantor_neighborhood":     guarantorNeighborhood,
		"guarantor_zip_code":         guarantorZipCode,
		"guarantor_city":             guarantorCity,
		"guarantor_birth_date":       guarantorBirthDate,
		"guarantor_uf":               guarantorUf,
		"guarantor_telephone_number": guarantorTelephoneNumber,
		"guarantor_address_code":     guarantorAddressCode,
	}

	clientInfo := map[string]interface{}{
		"cpf":              clientCpf,
		"name":             clientName,
		"email":            clientEmail,
		"cellphone_number": clientCellphoneNumber,
		"telephone_number": clientTelephoneNumber,
		"documents":        clientDocuments,
		"id":               clientId,
		"birth_date":       clientBirthDate,
	}

	enterpriseValueField := proposalEntity.mountIndividualField(databasePrefix, "property_value", proposalMap)
	if enterpriseValueField != nil {
		valueConverted, err := strconv.ParseFloat(enterpriseValueField.(string), 64)

		if err != nil {
			return nil, helper.ErrorBuilder(helper.ItWasNotPossibleConvertProprietyValueConst)
		}
		enterpriseValue = valueConverted
	}

	loanValue, err := strconv.ParseFloat(proposalMap[0]["proposal_down_payment"].(string), 64)
	if err != nil {
		return nil, helper.ErrorBuilder(helper.ItWasNotPossibleConvertLoanValueConst)
	}
	firstInstallmentValue, err := strconv.ParseFloat(proposalMap[0]["proposal_first_installment_value"].(string), 64)
	if err != nil {
		return nil, helper.ErrorBuilder(helper.ItWasNotPossibleConvertFirstInstallmentValueConst)
	}

	if proposalMap[0]["proposal_maximum_installment_value"] != nil {
		maximumInstallmentValue, err = strconv.ParseFloat(proposalMap[0]["proposal_maximum_installment_value"].(string), 64)
		if err != nil {
			return nil, helper.ErrorBuilder(helper.ItWasNotPossibleConvertMaximumInstallmentValueConst)
		}
	} else {
		maximumInstallmentValue = 0
	}

	if proposalMap[0]["proposal_approved"] != nil {
		proposalApproved = proposalMap[0]["proposal_approved"].(int64) > 0
	}

	if proposalMap[0]["proposal_formalized"] != nil {
		proposalFormalized = proposalMap[0]["proposal_formalized"].(int64) > 0
	}

	proposalResponse := map[string]interface{}{
		"proposal_id":             proposalMap[0]["proposal_id"],
		"proposal_status":         proposalMap[0]["proposal_status"],
		"proposal_number":         proposalMap[0]["proposal_number"],
		"operation":               proposalMap[0]["operation"],
		"company_down_payment":    proposalMap[0]["company_down_payment"],
		"difference_down_payment": proposalMap[0]["difference_down_payment"],
		"client":                  clientInfo,
		"guarantor":               guarantorInfo,
		"vehicle":                 proposalEntity.retrieveVehicleInfoAndOwner(proposalMap),
		"loan": proposalEntity.mountProposalLoanInfo(
			databasePrefix,
			proposalMap,
			loanValue,
			firstInstallmentValue,
			maximumInstallmentValue,
			proposalApproved,
			proposalFormalized,
			guarantorId,
		),
		"enterprise": proposalEntity.mountEnterpriseInfo(databasePrefix, proposalMap, enterpriseValue),
		"property":   proposalEntity.mountPropertyInfo(databasePrefix, proposalMap),
		"contract":   proposalEntity.mountContractAndSellerInfos(databasePrefix, proposalMap),
		"created_by": proposalEntity.mountCreatedByInfo(proposalMap),
		"messages": map[string]interface{}{
			"emcash":  proposalMap[0]["proposal_emcash_message"],
			"company": proposalMap[0]["proposal_company_message"],
		},
		"last_edit_user":      proposalEntity.mountLastEditUserInfo(proposalMap),
		"new_proposal":        proposalEntity.mountNewProposalInfo(proposalMap),
		"difference_proposal": proposalEntity.mountDifferenceProposalInfo(proposalMap),
		"base_proposal":       proposalEntity.mountBaseProposalInfo(proposalMap),
	}

	return proposalResponse, nil
}

func (proposalEntity *ProposalEntity) mountProposalLoanInfo(
	prefix string,
	proposal []map[string]interface{},
	loanValue float64,
	firstInstallmentValue float64,
	maximumInstallmentValue float64,
	proposalApproved bool,
	proposalFormalized bool,
	guarantorId string,
) map[string]interface{} {
	return map[string]interface{}{
		"value":                        loanValue,
		"buy_date":                     proposal[0]["proposal_buy_date"],
		"time":                         proposal[0]["proposal_installments"],
		"amortization_type":            proposal[0]["proposal_amortization_type"],
		"first_installment_value":      firstInstallmentValue,
		"maximum_installment_value":    maximumInstallmentValue,
		"payment_day":                  proposal[0]["proposal_payment_day"],
		"has_guarantor":                proposal[0]["proposal_has_guarantor"].(int64) > 0,
		"has_guarantee":                proposal[0]["proposal_has_guarantee"].(int64) > 0,
		"proposal_id":                  proposal[0]["proposal_id"],
		"emcash_risk_analysis_result":  proposal[0]["emcash_risk_analysis_result"],
		"company_risk_analysis_result": proposal[0]["company_risk_analysis_result"],
		"approved":                     proposalApproved,
		"formalized":                   proposalFormalized,
		"estimated_property_value":     proposalEntity.mountIndividualField(prefix, "proposal_estimated_property_value", proposal),
		"is_guarantor_refused":         proposal[0]["proposal_refused_guarantor"].(int64) > 0,
		"is_guarantee_refused":         proposal[0]["proposal_refused_guarantee"].(int64) > 0,
	}
}

func (proposalEntity *ProposalEntity) mountPropertyInfo(
	prefix string,
	proposalMap []map[string]interface{},
) map[string]interface{} {
	return map[string]interface{}{
		"value":            proposalEntity.mountIndividualField(prefix, "property_value", proposalMap),
		"complement":       proposalEntity.mountIndividualField(prefix, "property_complement", proposalMap),
		"apartment_number": proposalEntity.mountIndividualField(prefix, "property_apartment_number", proposalMap),
		"apartment_block":  proposalEntity.mountIndividualField(prefix, "property_apartment_block", proposalMap),
	}
}

func (proposalEntity *ProposalEntity) mountEnterpriseInfo(
	prefix string,
	proposal []map[string]interface{},
	enterpriseValue float64,
) map[string]interface{} {
	return map[string]interface{}{
		"id":               proposalEntity.mountIndividualField(prefix, "enterprise_id", proposal),
		"name":             proposalEntity.mountIndividualField(prefix, "enterprise_name", proposal),
		"value":            enterpriseValue,
		"complete_address": proposalEntity.mountIndividualField(prefix, "enterprise_complete_address", proposal),
		"number":           proposalEntity.mountIndividualField(prefix, "enterprise_address_code", proposal),
		"complement":       proposalEntity.mountIndividualField(prefix, "property_complement", proposal),
		"neighborhood":     proposalEntity.mountIndividualField(prefix, "enterprise_neighborhood", proposal),
		"zip_code":         proposalEntity.mountIndividualField(prefix, "enterprise_zip_code", proposal),
		"registry_office":  proposalEntity.mountIndividualField(prefix, "enterprise_registry_office", proposal),
		"uf":               proposalEntity.mountIndividualField(prefix, "enterprise_uf", proposal),
		"city":             proposalEntity.mountIndividualField(prefix, "enterprise_city", proposal),
	}
}

func (proposalEntity *ProposalEntity) mountCreatedByInfo(proposal []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"created_at": proposal[0]["proposal_created_at"],
		"user_name":  proposal[0]["creator_name"],
		"user_email": proposal[0]["creator_email"],
		"user_id":    proposal[0]["creator_id"],
	}
}

func (proposalEntity *ProposalEntity) mountLastEditUserInfo(proposal []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"updated_at": proposal[0]["proposal_updated_at"],
		"user_name":  proposal[0]["proposal_last_edit_name"],
		"user_email": proposal[0]["proposal_last_edit_email"],
		"company":    "MRV",
		"branch":     "N/A",
		"product":    "N/A",
	}
}

func (proposalEntity *ProposalEntity) mountNewProposalInfo(proposal []map[string]interface{}) map[string]interface{} {
	var (
		hasGuarantor bool
		hasGuarantee bool
	)

	if proposal[0]["proposal_has_guarantor"] != nil {
		hasGuarantor = proposal[0]["proposal_has_guarantor"].(int64) > 0
	}

	if proposal[0]["proposal_has_guarantee"] != nil {
		hasGuarantee = proposal[0]["proposal_has_guarantee"].(int64) > 0
	}

	if proposal[0]["proposal_status"] == proposalDomain.ProposalStatusRefusedConst {
		return nil
	}

	return map[string]interface{}{
		"new_proposal_id":                        proposal[0]["new_proposal_id"],
		"new_proposal_installments":              proposal[0]["new_proposal_installments"],
		"new_proposal_created_at":                proposal[0]["new_proposal_created_at"],
		"new_proposal_has_guarantor":             hasGuarantor,
		"new_proposal_has_guarantee":             hasGuarantee,
		"new_proposal_down_payment":              proposal[0]["new_proposal_down_payment"],
		"new_proposal_payment_day":               proposal[0]["new_proposal_payment_day"],
		"new_emcash_risk_analysis_result":        proposal[0]["new_emcash_risk_analysis_result"],
		"new_company_risk_analysis_result":       proposal[0]["new_company_risk_analysis_result"],
		"new_proposal_amortization_type":         proposal[0]["new_proposal_amortization_type"],
		"new_proposal_buy_date":                  proposal[0]["new_proposal_buy_date"],
		"new_proposal_difference_down_payment":   proposal[0]["new_proposal_difference_down_payment"],
		"new_proposal_company_down_payment":      proposal[0]["new_proposal_company_down_payment"],
		"new_proposal_maximum_installment_value": proposal[0]["new_proposal_maximum_installment_value"],
	}
}

func (proposalEntity *ProposalEntity) mountContractAndSellerInfos(
	prefix string,
	proposal []map[string]interface{},
) map[string]interface{} {
	var (
		contractId            string
		contractDate          string
		contractNumber        string
		contractMatriculation string
		sellerBusinessName    string
		sellerCnpj            string
		sellerId              string
	)

	contractIdField := proposalEntity.mountIndividualField(prefix, "contract_id", proposal)
	if contractIdField != nil {
		contractId = contractIdField.(string)
	}

	contractDateField := proposalEntity.mountIndividualField(prefix, "contract_date", proposal)
	if contractDateField != nil {
		contractDate = contractDateField.(string)
	}

	contractNumberField := proposalEntity.mountIndividualField(prefix, "contract_number", proposal)
	if contractNumberField != nil {
		contractNumber = contractNumberField.(string)
	}

	contractMatriculationField := proposalEntity.mountIndividualField(prefix, "contract_matriculation", proposal)
	if contractMatriculationField != nil {
		contractMatriculation = contractMatriculationField.(string)
	}

	sellerBusinessNameField := proposalEntity.mountIndividualField(prefix, "seller_business_name", proposal)
	if sellerBusinessNameField != nil {
		sellerBusinessName = sellerBusinessNameField.(string)
	}

	sellerCnpjField := proposalEntity.mountIndividualField(prefix, "seller_cnpj", proposal)
	if sellerCnpjField != nil {
		sellerCnpj = sellerCnpjField.(string)
	}

	sellerIdField := proposalEntity.mountIndividualField(prefix, "seller_id", proposal)
	if sellerIdField != nil {
		sellerId = sellerIdField.(string)
	}

	return map[string]interface{}{
		"contract_id":            contractId,
		"contract_date":          contractDate,
		"contract_number":        contractNumber,
		"contract_matriculation": contractMatriculation,
		"seller_business_name":   sellerBusinessName,
		"seller_cnpj":            sellerCnpj,
		"seller_id":              sellerId,
	}
}

func (proposalEntity *ProposalEntity) retrieveVehicleInfoAndOwner(proposal []map[string]interface{}) map[string]interface{} {
	for _, rangeProposal := range proposal {
		if rangeProposal["vehicle_owner_cpf"] != nil {
			return map[string]interface{}{
				"brand":   rangeProposal["vehicle_brand"],
				"model":   rangeProposal["vehicle_model"],
				"year":    rangeProposal["vehicle_year"],
				"plate":   rangeProposal["vehicle_plate"],
				"renavam": rangeProposal["vehicle_renavam"],
				"value":   rangeProposal["vehicle_guarantee_value"],
				"uf":      rangeProposal["vehicle_uf"],
				"city":    rangeProposal["vehicle_city"],
				"color":   rangeProposal["vehicle_color"],
				"owner": map[string]interface{}{
					"type": rangeProposal["person_type"].(string),
					"cpf":  rangeProposal["vehicle_owner_cpf"],
					"id":   rangeProposal["vehicle_owner_id"],
				},
			}
		}
	}

	return map[string]interface{}{
		"brand":   nil,
		"model":   nil,
		"year":    nil,
		"plate":   nil,
		"renavam": nil,
		"value":   nil,
		"uf":      nil,
		"city":    nil,
		"color":   nil,
		"owner": map[string]interface{}{
			"type": nil,
			"cpf":  nil,
			"id":   nil,
		},
	}
}

func (proposalEntity *ProposalEntity) VerifyIfProposalCanBeApprovedOrFormalized(
	proposal map[string]interface{},
	proposalApprovedAndFormalizedInfo *proposalDomain.ProposalApprovedAndFormalized,
) error {
	proposalApproved := proposal["loan"].(map[string]interface{})["approved"].(bool)
	proposalFormalized := proposal["loan"].(map[string]interface{})["formalized"].(bool)
	hasGuarantor := proposal["loan"].(map[string]interface{})["has_guarantor"].(bool)
	hasGuarantee := proposal["loan"].(map[string]interface{})["has_guarantee"].(bool)
	proposalStatus := proposal["proposal_status"].(string)

	if proposalApproved && !proposalApprovedAndFormalizedInfo.Approved {
		return helper.ErrorBuilder(helper.ProposalCannotBeUnapprovedConst)
	}

	if proposalFormalized && !proposalApprovedAndFormalizedInfo.Formalized {
		return helper.ErrorBuilder(helper.ProposalCannotBeDeformalizedConst)
	}

	if !proposalApproved && proposalApprovedAndFormalizedInfo.Approved &&
		proposalStatus != proposalDomain.ProposalStatusApprovedConst {
		return helper.ErrorBuilder(helper.ProposalMustBeIntoStatusApprovedToBeSetAsWaitingSignatureConst)
	}

	isOnPreContractIssuedOrContractSettled := proposalStatus == proposalDomain.ProposalStatusPreContractIssuedConst || proposalStatus == proposalDomain.ProposalStatusContractSettledConst
	shouldFormalizeProposal := !proposalFormalized && proposalApprovedAndFormalizedInfo.Formalized

	if shouldFormalizeProposal && !isOnPreContractIssuedOrContractSettled {
		return helper.ErrorBuilder(helper.ProposalMustBeInPreContractIssuedToBeSetAsFormalizedConst)
	}

	for key, value := range proposal {
		if value == nil {
			continue
		}

		isMap := reflect.TypeOf(value)

		switch key {
		case proposalDomain.GuarantorConst:
			if !hasGuarantor {
				continue
			}
		case proposalDomain.VehicleConst:
			if !hasGuarantee {
				continue
			}
		case proposalDomain.CreatedByConst,
			proposalDomain.LastEditUserConst,
			proposalDomain.NewProposalConst,
			proposalDomain.MessagesConst,
			proposalDomain.BaseProposalConst,
			proposalDomain.DifferenceProposalConst:
			continue
		}

		if isMap.Kind() == reflect.Map && isMap.Key().Kind() == reflect.String && isMap.Elem().Kind() == reflect.Interface {
			for innerKey, innerValue := range value.(map[string]interface{}) {
				if innerValue == nil {
					return helper.ErrorBuilder(helper.ProposalHasMandatoryInformationNotFilledInConst)
				}

				isInterMap := reflect.TypeOf(innerValue)
				if isInterMap.Kind() == reflect.Map && isInterMap.Key().Kind() == reflect.String && isInterMap.Elem().Kind() == reflect.Interface {
					for deepKey, deepValue := range innerValue.(map[string]interface{}) {
						if deepValue == nil {
							return helper.ErrorBuilder(helper.ProposalHasMandatoryInformationNotFilledInConst)
						}

						if err := proposalEntity.throwErrorIfValueEmpty(deepKey, deepValue, proposalStatus); err != nil {
							return err
						}
					}
				}

				if innerKey == contractDomain.ContractNumberConst && proposalApprovedAndFormalizedInfo.Approved {
					continue
				}

				if err := proposalEntity.throwErrorIfValueEmpty(innerKey, innerValue, proposalStatus); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (proposalEntity *ProposalEntity) throwErrorIfValueEmpty(key string, value interface{}, proposalStatus string) error {
	var mandatoryKeys = map[string]bool{
		"has_guarantee":                true,
		"proposal_id":                  true,
		"emcash_risk_analysis_result":  true,
		"formalized":                   true,
		"value":                        true,
		"buy_date":                     true,
		"amortization_type":            true,
		"company_risk_analysis_result": true,
		"time":                         true,
		"first_installment_value":      true,
		"payment_day":                  true,
		"cellphone_number":             true,
		"documents":                    true,
		"birth_date":                   true,
		"cpf":                          true,
		"name":                         true,
		"email":                        true,
		"apartment_number":             true,
		"apartment_block":              true,
		"complement_address":           true,
		"neighborhood":                 true,
		"zip_code":                     true,
		"uf":                           true,
		"city":                         true,
		"plate":                        true,
		"renavam":                      true,
		"owner":                        true,
		"brand":                        true,
		"year":                         true,
		"color":                        true,
		"contract_matriculation":       true,
		"seller_business_name":         true,
		"seller_cnpj":                  true,
		"contract_date":                true,
		"contract_number":              true,
	}

	if proposalStatus == proposalDomain.ProposalStatusPreContractIssuedConst {
		mandatoryKeys["contract_id"] = true
	}

	if mandatoryKeys[key] {
		switch val := value.(type) {
		case string:
			if val == "" {
				return helper.ErrorBuilder(helper.ProposalHasMandatoryInformationNotFilledInConst)
			}
		case float64:
			if val == 0 {
				return helper.ErrorBuilder(helper.ProposalHasMandatoryInformationNotFilledInConst)
			}
		case int64:
			if val == 0 {
				return helper.ErrorBuilder(helper.ProposalHasMandatoryInformationNotFilledInConst)
			}
		case []map[string]interface{}:
			if key == proposalDomain.DocumentsConst && len(val) < 4 {
				return helper.ErrorBuilder(helper.ProposalHasMandatoryInformationNotFilledInConst)
			}
		}
	}

	return nil
}

func (proposalEntity *ProposalEntity) mountDifferenceProposalInfo(proposal []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"difference_id":                           proposal[0]["difference_proposal_id"],
		"difference_proposal_number":              proposal[0]["difference_proposal_number"],
		"difference_proposal_status":              proposal[0]["difference_proposal_status"],
		"difference_installments":                 proposal[0]["difference_proposal_installments"],
		"difference_created_at":                   proposal[0]["difference_proposal_created_at"],
		"difference_has_guarantor":                false,
		"difference_has_guarantee":                false,
		"difference_down_payment":                 proposal[0]["difference_proposal_down_payment"],
		"difference_payment_day":                  proposal[0]["difference_proposal_payment_day"],
		"difference_emcash_risk_analysis_result":  proposal[0]["difference_proposal_emcash_risk_analysis_result"],
		"difference_company_risk_analysis_result": proposal[0]["difference_proposal_company_risk_analysis_result"],
		"difference_amortization_type":            proposal[0]["difference_proposal_amortization_type"],
		"difference_buy_date":                     proposal[0]["difference_proposal_buy_date"],
	}
}

func (proposalEntity *ProposalEntity) mountBaseProposalInfo(proposal []map[string]interface{}) map[string]interface{} {
	var (
		hasGuarantor bool
		hasGuarantee bool
	)

	if proposal[0]["base_proposal_has_guarantor"] != nil {
		hasGuarantor = proposal[0]["base_proposal_has_guarantor"].(int64) > 0
	}

	if proposal[0]["base_proposal_has_guarantee"] != nil {
		hasGuarantee = proposal[0]["base_proposal_has_guarantee"].(int64) > 0
	}

	return map[string]interface{}{
		"base_proposal_id":                  proposal[0]["base_proposal_id"],
		"base_proposal_number":              proposal[0]["base_proposal_number"],
		"base_proposal_status":              proposal[0]["base_proposal_status"],
		"base_proposal_installments":        proposal[0]["base_proposal_installments"],
		"base_proposal_created_at":          proposal[0]["base_proposal_created_at"],
		"base_proposal_has_guarantor":       hasGuarantor,
		"base_proposal_has_guarantee":       hasGuarantee,
		"base_proposal_down_payment":        proposal[0]["base_proposal_down_payment"],
		"base_proposal_payment_day":         proposal[0]["base_proposal_payment_day"],
		"base_emcash_risk_analysis_result":  proposal[0]["base_proposal_emcash_risk_analysis_result"],
		"base_company_risk_analysis_result": proposal[0]["base_proposal_company_risk_analysis_result"],
		"base_proposal_amortization_type":   proposal[0]["base_proposal_amortization_type"],
		"base_proposal_buy_date":            proposal[0]["base_proposal_buy_date"],
	}
}

func (proposalEntity *ProposalEntity) mountIndividualField(
	prefix string,
	sufix string,
	proposal []map[string]interface{},
) any {
	return proposal[0][prefix+sufix]
}
