package core

type Loc struct {
	Filename   string
	Line, Char int
}

func MakeLoc(Filename string, Line, Char int) *Loc {
	return &Loc{Filename, Line, Char}
}

func (l *Loc) Error(e error) *PhpError {
	switch err := e.(type) {
	case *PhpError:
		if err.l == nil {
			err.l = l
		}
		return err
	default:
		return &PhpError{e, l, PhpErrorFatal}
	}
}