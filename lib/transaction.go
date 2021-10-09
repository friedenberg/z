package lib

type Transaction struct {
	Add *transactionPrinter
	Mod *transactionPrinter
	Del *transactionPrinter
}

type ZettelSlice []*Zettel

func (s ZettelSlice) Paths() (p []string) {
	p = make([]string, len(s))

	for i, z := range s {
		p[i] = z.Path
	}

	return
}

func (t Transaction) Added() ZettelSlice {
	return t.Add.zettels
}

func (t Transaction) Modified() ZettelSlice {
	return t.Mod.zettels
}

func (t Transaction) Deleted() ZettelSlice {
	return t.Del.zettels
}

func (t Transaction) RawFiles() (f []string) {
	f = make([]string, 0, len(t.Added())+len(t.Modified())+len(t.Deleted()))

	add := func(s []*Zettel) {
		for _, z := range s {
			f = append(f, z.Path)
		}
	}

	add(t.Added())
	add(t.Modified())
	add(t.Deleted())

	return
}

type transactionPrinter struct {
	zettels []*Zettel
	files   []string
	urls    []string
	errors  []error
}

func (p *transactionPrinter) Begin() {}
func (p *transactionPrinter) End()   {}
func (p *transactionPrinter) PrintZettel(i int, z *Zettel, err error) {
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
