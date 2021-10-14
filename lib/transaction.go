package lib

func MakeTransaction() (t Transaction) {
	t = Transaction{
		Add: &transactionPrinter{},
		Mod: &transactionPrinter{},
		Del: &transactionPrinter{},
	}

	return
}

//TODO-P1 handle cases where files or zettles are just opened but not edited
type Transaction struct {
	ShouldSkipCommit bool
	Add              *transactionPrinter
	Mod              *transactionPrinter
	Del              *transactionPrinter
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

	if f, ok := z.Note.Metadata.LocalFile(); ok {
		//TODO-P2 support deleting local files
		p.files = append(p.files, f.Id)
	}

	if u, ok := z.Note.Metadata.Url(); ok {
		p.urls = append(p.urls, u.String())
	}

	return
}
