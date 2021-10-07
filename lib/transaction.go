package lib

type Transaction struct {
	Add transactionPrinter
	Mod transactionPrinter
	Del transactionPrinter
}

func (t Transaction) Added() []*Zettel {
	return t.Add.zettels
}

func (t Transaction) Modified() []*Zettel {
	return t.Mod.zettels
}

func (t Transaction) Deleted() []*Zettel {
	return t.Del.zettels
}

type transactionPrinter struct {
	zettels []*Zettel
	files   []string
	urls    []string
	errors  []error
}

func (p transactionPrinter) Begin() {}
func (p transactionPrinter) End()   {}
func (p transactionPrinter) PrintZettel(i int, z *Zettel, err error) {
	if err != nil {
		p.errors = append(p.errors, err)
		return
	}

	p.zettels = append(p.zettels, z)

	if z.HasFile() {
		//TODO which filepath?
		p.files = append(p.files, z.FilePath())
	}

	if z.HasUrl() {
		p.urls = append(p.urls, z.Metadata.Url)
	}

	return
}
