package text

import "strings"

type Text string

func (t Text) String() string {
	return string(t)
}

func (t Text) ToLower() Text {
	return Text(strings.ToLower(t.String()))
}

func (t Text) ToUpper() Text {
	return Text(strings.ToUpper(t.String()))
}

func (t Text) ToSnake() Text {
	panic("not implemented")
}

func (t Text) ToScreamingSnake() Text {
	// return Text(strcase.ToScreamingSnake(t.String()))
	panic("not implemented")
}

func (t Text) ReplaceAll(old string, new string) Text {
	return Text(strings.Replace(t.String(), old, new, -1))
}

func (t Text) Replace(old, new string, n int) Text {
	return Text(strings.Replace(t.String(), old, new, n))
}
