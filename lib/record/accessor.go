package record

import "errors"

type Accessor struct {
	byIndex      bool
	bySimpleName bool
	index        int
	name         string
}

func (a *Accessor) getValue(rec record) (string, error) {
	switch {
	case a.byIndex:
		valuerec, err := rec.GetByIndex(a.index)
		if err != nil {
			return "", nil
		}

		value, err := valuerec.AsValue()
		if err != nil {
			return "", err
		}

		return value, nil

	case a.bySimpleName:
		valuerec, err := rec.GetByName(a.name)
		if err != nil {
			return "", nil
		}

		value, err := valuerec.AsValue()
		if err != nil {
			return "", err
		}

		return value, nil

	default:
		return "", errors.New("invalid accessor")
	}
}
