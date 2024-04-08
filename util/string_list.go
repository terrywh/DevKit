package util

type StringList []string

func (i *StringList) String() string {
	return "my string representation"
}

func (i *StringList) Set(value string) error {
	*i = append(*i, value)
	return nil
}
