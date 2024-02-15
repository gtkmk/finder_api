package port

type FileFactoryInterface interface {
	Make(string) FileInterface
}
