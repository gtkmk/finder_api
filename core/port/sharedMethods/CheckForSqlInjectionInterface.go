package sharedMethods

type CheckForSqlInjectionInterface interface {
	CheckForSqlInjection(input string) error
}
