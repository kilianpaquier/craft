package generate

// Delimiter represents the pair of start and end delimiter for go template substitution.
type Delimiter struct {
	// EndDelim is the end delimiter of a go template statement, i.e. >> or }} or ]], etc.
	EndDelim string

	// StartDelim is the start delimiter of a go template statement, i.e. << or {{ or [[, etc.
	StartDelim string
}

var (
	chevron = Delimiter{
		EndDelim:   ">>",
		StartDelim: "<<",
	}

	bracket = Delimiter{
		EndDelim:   "}}",
		StartDelim: "{{",
	}

	squareBracket = Delimiter{
		EndDelim:   "]]",
		StartDelim: "[[",
	}
)

// DelimiterChevron returns go template delimiter << and >>.
func DelimiterChevron() Delimiter {
	return chevron
}

// DelimiterBracket returns go template delimiter {{ and }}.
func DelimiterBracket() Delimiter {
	return bracket
}

// DelimiterSquareBracket returns go template delimiter [[ and ]].
func DelimiterSquareBracket() Delimiter {
	return squareBracket
}
