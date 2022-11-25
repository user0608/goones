package answer

type Pager interface {
	Data() any
	ReturnedItems() int
	RequestedItems() int
	CurrentPage() int
	NumberPages() int
	NumberItems() int
}
