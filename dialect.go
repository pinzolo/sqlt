package sqlt

// Dialect resolves dialect of each databse.
type Dialect interface {
	// IsOrdinalPlaceholderSupportedreturns true if databse support ordinal placeholder.
	// ex: $1, $2 (PostgreSQL)
	IsOrdinalPlaceholderSupported() bool
	// OrdinalPlaceholderPrefix returns prefix of placeholder.
	// This is used when IsOrdinalPlaceholderSupported is true in Exec func.
	OrdinalPlaceHolderPrefix() string
	// Placeholder character.
	// This is used when IsOrdinalPlaceholderSupported is false.
	Placeholder() string
	// NamedPlaceholderPrefix returns prefix of placeholder.
	// This is used in ExecNamed func.
	NamedPlaceholderPrefix() string
}

var (
	// Postgres is PostgreSQL dialect resolver.
	Postgres = postgres{}
	// MySQL is MySQL dialect resolver.
	MySQL = mysql{}
)

type postgres struct{}

func (p postgres) IsOrdinalPlaceholderSupported() bool {
	return true
}

func (p postgres) OrdinalPlaceHolderPrefix() string {
	return "$"
}

func (p postgres) Placeholder() string {
	return ""
}

func (p postgres) NamedPlaceholderPrefix() string {
	return ":"
}

type mysql struct{}

func (m mysql) IsOrdinalPlaceholderSupported() bool {
	return false
}

func (m mysql) OrdinalPlaceHolderPrefix() string {
	return ""
}

func (m mysql) Placeholder() string {
	return "?"
}

func (m mysql) NamedPlaceholderPrefix() string {
	return ":"
}