package forms

type InfoSwaggerForm struct {
	ProductId   string
	ProductName string
	Quantity    int
	Price       int
}
type OrderSwaggerForm struct {
	Total    int
	Tax      float32
	Currency string
	Info     *[]InfoSwaggerForm
}
