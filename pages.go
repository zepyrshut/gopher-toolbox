package gorender

// Pages contiene la información de paginación.
type Pages struct {
	// totalElements son la cantidad de elementos totales a paginar. Pueden ser
	// total de filas o total de páginas de blog.
	totalElements int
	// showElements muestra la cantidad máxima de elementos a mostrar en una
	// página.
	showElements int
	// currentPage es la página actual, utilizado como ayuda para mostrar la
	// página activa.
	currentPage int
}

// Page contiene la información de una página.
type Page struct {
	// number es el número de página.
	number int
	// active es un dato lógico que indica si la página es la actual.
	active bool
}

// NewPages crea un nuevo objeto para paginación.
func NewPages(totalElements, showElements, currentPage int) Pages {
	if showElements <= 0 {
		showElements = 1
	}
	if currentPage <= 0 {
		currentPage = 1
	}
	p := Pages{totalElements, showElements, currentPage}
	if p.currentPage > p.TotalPages() {
		p.currentPage = p.TotalPages()
	}

	return p
}

// Limit devuelve la cantidad de elementos máximos a mostrar por página.
func (p *Pages) Limit() int {
	return p.showElements
}

// TotalPages devuelve la cantidad total de páginas.
func (p *Pages) TotalPages() int {
	return (p.totalElements + p.showElements - 1) / p.showElements
}

// IsFirst indica si la página actual es la primera.
func (p *Pages) IsFirst() bool {
	return p.currentPage == 1
}

// IsLast indica si la página actual es la última.
func (p *Pages) IsLast() bool {
	return p.currentPage == p.TotalPages()
}

// HasPrevious indica si hay una página anterior.
func (p *Pages) HasPrevious() bool {
	return p.currentPage > 1
}

// HasNext indica si hay una página siguiente.
func (p *Pages) HasNext() bool {
	return p.currentPage < p.TotalPages()
}

// Previous devuelve el número de la página anterior.
func (p *Pages) Previous() int {
	return p.currentPage - 1
}

// Next devuelve el número de la página siguiente.
func (p *Pages) Next() int {
	return p.currentPage + 1
}

func (p *Page) NumberOfPage() int {
	return p.number
}

// IsActive indica si la página es la actual.
func (p *Page) IsActive() bool {
	return p.active
}

// Pages devuelve un arreglo de páginas para mostrar en la paginación. El
// parametro pagesShow indica la cantidad de páginas a mostrar, asignable desde
// la plantilla.
func (p *Pages) Pages(pagesShow int) []*Page {
	var pages []*Page
	startPage := p.currentPage - (pagesShow / 2)
	endPage := p.currentPage + (pagesShow/2 - 1)

	if startPage < 1 {
		startPage = 1
		endPage = pagesShow
	}

	if endPage > p.TotalPages() {
		endPage = p.TotalPages()
		startPage = p.TotalPages() - pagesShow + 1
		if startPage < 1 {
			startPage = 1
		}
	}

	for i := startPage; i <= endPage; i++ {
		pages = append(pages, &Page{i, i == p.currentPage})
	}
	return pages
}

func PaginateArray[T any](items []T, currentPage, itemsPerPage int) []T {
	totalItems := len(items)

	startIndex := (currentPage - 1) * itemsPerPage
	endIndex := startIndex + itemsPerPage

	if startIndex > totalItems {
		startIndex = totalItems
	}
	if endIndex > totalItems {
		endIndex = totalItems
	}

	return items[startIndex:endIndex]
}
