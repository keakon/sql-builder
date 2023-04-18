package sb

type AliasMode uint8

const (
	NoAlias AliasMode = iota
	UseAlias
	OnlyAlias
)
