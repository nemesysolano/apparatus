package embeddings

import "errors"

type InstitutionCode int

const (
	DGII InstitutionCode = iota
	SB
	SIMV
)

var InstitutionAcronyms = map[string]InstitutionCode{
	"DGII": DGII,
	"SB":   SB,
	"SIMV": SIMV,
}

func (i InstitutionCode) String() (string, error) {
	for acronym, code := range InstitutionAcronyms {
		if code == i {
			return acronym, nil
		}
	}
	return "", errors.New("invalid institution code")
}
