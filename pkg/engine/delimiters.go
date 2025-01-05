package engine

// Delimiters represents the pair of start and end delimiter for go template substitution.
type Delimiters struct {
	// EndDelim is the end delimiter of a go template statement, i.e. >> or }} or ]], etc.
	EndDelim string

	// StartDelim is the start delimiter of a go template statement, i.e. << or {{ or [[, etc.
	StartDelim string
}

var (
	chevron = Delimiters{
		EndDelim:   ">>",
		StartDelim: "<<",
	}

	bracket = Delimiters{
		EndDelim:   "}}",
		StartDelim: "{{",
	}

	squareBracket = Delimiters{
		EndDelim:   "]]",
		StartDelim: "[[",
	}
)

// DelimitersChevron returns go template delimiter << and >>.
func DelimitersChevron() Delimiters {
	return chevron
}

// DelimitersBracket returns go template delimiter {{ and }}.
func DelimitersBracket() Delimiters {
	return bracket
}

// DelimitersSquareBracket returns go template delimiter [[ and ]].
func DelimitersSquareBracket() Delimiters {
	return squareBracket
}
