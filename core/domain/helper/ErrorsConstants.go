package helper

import (
	"fmt"
	"strings"
)

const hasParametersSymbol = "%"

func ErrorBuilder(messageType string, params ...any) error {
	msg := DefineMassage(messageType, params...)

	return fmt.Errorf(msg)
}

func DefineMassage(msg string, params ...any) string {
	if defineShouldHasParams(msg) {
		return fmt.Sprintf(msg, params...)
	}

	return msg
}

func defineShouldHasParams(msg string) bool {
	return strings.Contains(msg, hasParametersSymbol)
}

const (
	SendAValidCNPJErrorMessageConst                      = "Envie um CNPJ válido"
	InvalidCnpjWithOnlyZerosMessageConst                 = "CNPJ inválido."
	InvalidCnpjLengthMessageConst                        = "O CNPJ não possui o tamanho correto."
	CnpjMustHaveOnlyNumberMessageConst                   = "O CNPJ deve conter apenas números."
	PasswordMustHaveCharactersMessageConst               = "A senha deve conter letras e números."
	PasswordMustHaveAtLastOneUpperCaseMessageConst       = "A senha deve conter pelo menos uma letra maiúscula."
	PasswordMustHaveAtLastOneDigitMessageConst           = "A senha deve conter pelo menos um número."
	PasswordMustHaveAtLastOneSpecialCharacterConst       = "A senha deve conter pelo menos um caractere especial."
	PasswordMustHaveTheMinLengthMessageConst             = "A senha precisa ter ao mínimo %d caracteres"
	PasswordMayNotEmptyMessageConst                      = "A senha não pode ser vazia."
	IsNecessaryToSendTheOldPasswordMessageConst          = "É necessário o envio da senha antiga."
	PasswordsNotEqualMessageConst                        = "As senhas não coincidem."
	EmptyJsonDataMessageConst                            = "Dados do json vazio."
	JsonNotFoundMessageConst                             = "Json não encontrado."
	EmailOrPasswordIncorrectConst                        = "E-mail ou senha incorretos."
	SomethingWentWrongConst                              = "Algo de errado aconteceu= %s"
	UnauthorizedConst                                    = "Não autorizado"
	ExpiredTokenConst                                    = "Token expirado"
	ErrorParsingTokenConst                               = "Erro ao parsear token"
	ErrorGeneratingEncryptionConst                       = "Falha ao gerar criptografia= %s"
	UserIdIsRequiredConst                                = "O id do usuário é obrigatório."
	ErrorExtractingUserIdConst                           = "Falha ao extrair o id do usuário."
	ErrorExtractingUserGroupConst                        = "Falha ao extrair o grupo do usuário."
	UserNotFoundConst                                    = "Usuário não encontrado."
	UserAlreadyRegisteredConst                           = "Esse usuário já foi cadastrado."
	UserAlreadyRegisteredWithEmailConst                  = "Já existe um usuário cadastrado com esse e-mail."
	ContentTypeErrorConst                                = "Content-Type da requisição precisa ser multipart/form-data"
	NoRecordsFoundConst                                  = "Nenhum registro encontrado."
	InvalidEmailFormatConst                              = "O e-mail possui um formato inválido."
	EmailCannotBeEmptyConst                              = "E-mail não pode ser vazio."
	EmailTooLongConst                                    = "E-mail não pode ter mais de 50 caracteres."
	ErrorWithCodeConst                                   = "Erro interno= código %d."
	UnknownResultTypeConst                               = "Tipo de resultado desconhecido. Verifique default do método handleScanReturn."
	FunctionalityNotImplementedConst                     = "Funcionalidade ainda não implementeda."
	TimeoutLimitReachedConst                             = "Timeout"
	InvalidDateFormatConst                               = "Formato de data inválido."
	FieldCannotHaveMoreThanSetCharactersConst            = "%s não pode ter mais de %d caracteres."
	FieldCannotBeEmptyConst                              = "%s não pode ser vazio."
	PasswordCannotBeSameAsOldPasswordConst               = "A nova senha não pode ser igual a antiga."
	OldPasswordIncorrectConst                            = "Senha antiga incorreta."
	TokenExpiredGenerateNewEmailConst                    = "Convite expirado, gere um e-mail novamente."
	EmailSentToChangePasswordConst                       = "Um e-mail foi enviado para você trocar a senha."
	ErrorShuttingDownServerConst                         = "Error shutting down server= %v\n."
	ThePhoneShouldOnlyCountNumbersConst                  = "O número de telefone deve conter somente números."
	ThePhoneIsNotInAvalidLengthConst                     = "O número de telefone não possui o tamanho correto."
	UserDoNotRequestPasswordChangingConst                = "Este usuário não solicitou a troca de senha"
	IncorrectPasswordOrLoginConst                        = "Senha ou e-mail incorretos"
	InviteExpiredGenerateNewEmailConst                   = "Convite expirado, gere um e-mail novamente."
	EmailToChangePasswordSentConst                       = "Um e-mail foi enviado para você trocar a senha."
	TheGroupPermissionLayerCannotBeResetConst            = "A camada de permissão do grupo não pode ser zerada."
	TheGroupPermissionLayerCannotBeNegativeConst         = "A camada de permissão do grupo não pode ser negativa."
	TheGroupPermissionLayerCannotBeGreaterThanLimitConst = "A camada de permissão do grupo não pode ser maior que 999."
	PasswordsAreDifferentsConst                          = "As senhas não coincidem."
	FieldIsMandatoryConst                                = "O campo '%s' é obrigatório"
	FieldIsMandatoryAndMustToBeGreaterThanZeroConst      = "O campo '%s' é obrigatório e deve ser maior que zero."
	InvalidFileTypeConst                                 = "Tipo de arquivo inválido= %s"
	UserIsNotActiveConst                                 = "Houve um problema ao tentar acessar a plataforma. Por favor, entre em contato com o suporte."
	InformFieldConst                                     = "É necessário informar %s."
	IdIsNotAValidUuidConst                               = "O id não é um uuid válido."
	ErrorConvertingValueConst                            = "Erro ao converter valor para '%s'."
	ErrorInvalidFileTypeConst                            = "Tipo do arquivo %s inválido. Tente enviar em outro formato. Por exemplo= jpg, jpeg ou png"
	FirstAccessLinkAlreadyUsedConst                      = "Esse link de redefinição já foi utilizado! Gere outro link."
	PasswordMissingUppercaseConst                        = "A senha deve conter pelo menos uma letra maiúscula"
	PasswordMissingLowercaseConst                        = "A senha deve conter pelo menos uma letra minúscula"
	PasswordMissingNumberConst                           = "A senha deve conter pelo menos um número"
	PasswordMissingSpecialCharConst                      = "A senha deve conter pelo menos um caractere especial"
	ErrorGettingFormFileConst                            = "Erro ao obter arquivo do formulário: %s"
	ErrorWhenGettingFormFileWithoutInterpolationConst    = "Erro ao obter arquivo do formulário."
	FieldNotInAllowedValuesConst                         = "O valor para %s não é permitido."
	InvalidPageNumberErrorConst                          = "Número de página inválido."
	OptionNotRecognizedMessageConst                      = "Opção de '%s' não reconhecida."
	PostLostFoundStatusNotRecognizedConst                = "Status de achado ou perdido não reconhecido."
	PostRewardNotRecognizedMessageConst                  = "Opção de recompensa não reconhecida."
	PostNotFoundMessageConst                             = "Postagem não encontrada."
	CommentNotFoundMessageConst                          = "Comentário não encontrado."
	ErrorReadingFileConst                                = "Erro ao obter arquivo: %s"
	InvalidLikeTypeConst                                 = "Tipo de like inválido para o post ou o comentário."
	InvalidLikeRequestConst                              = "Requisição de like inválida. Definir destino do like, post ou comentário."
)

const (
	ErrorWhenTryToCommitTransactionCodeConst           = 1003
	ErrorWhenTryToRollbackTransactionCodeConst         = 1004
	ErrorWhenTryToCreateATransactionSavePointCodeConst = 1005
	ErrorWhenTryToRollbackToSavePointCodeConst         = 1006
	ErrorWhenExecuteRawStatementCodeConst              = 1007
	ErrorWhenTryToExecuteRowsQueryCodeConst            = 1008
	ErrorLoadingSaoPauloLocationCodeConst              = 2000
	ErrorGeneratingFormattedTimestampCodeConst         = 2001
)
