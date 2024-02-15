package repositories

type ErrorRepositoryInterface interface {
	SaveErrorWithoutStack(id string, errorMessage string, createdAt string)
	SaveErrorWithStack(id string, errorMessage string, stack string, createdAt string)
}
