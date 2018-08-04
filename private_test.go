package sqlt

import "testing"

func TestDropSample(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2, 3)
	AND created_at = /*%p "time"%*/'2000-01-01 12:34:56'
	AND name = /*%p "name"%*/'foo'
	AND age = /*%p "age"%*/18`
	expected := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/
	AND created_at = /*%p "time"%*/
	AND name = /*%p "name"%*/
	AND age = /*%p "age"%*/`
	if actual := dropSample(s); actual != expected {
		t.Errorf("dropSample faild: expected %s, but got %s", expected, actual)
	}
}
