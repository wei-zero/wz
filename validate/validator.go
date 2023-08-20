package validate

type Validator interface {
	Validate() error
}

func Validate(v any) error {
	if vx, ok := v.(Validator); ok {
		return vx.Validate()
	}

	return nil
}
