package log

import "go.uber.org/zap"

type Field = zap.Field

var (
	Any      = zap.Any
	Binary   = zap.Binary
	Bool     = zap.Bool
	Float32  = zap.Float32
	Float64  = zap.Float64
	Int      = zap.Int
	Int8     = zap.Int8
	Int16    = zap.Int16
	Int32    = zap.Int32
	Int64    = zap.Int64
	String   = zap.String
	Stringer = zap.Stringer
	Uint     = zap.Uint
	Uint8    = zap.Uint8
	Uint16   = zap.Uint16
	Uint32   = zap.Uint32
	Uint64   = zap.Uint64
	StdError = zap.Error
	Time     = zap.Time
	Duration = zap.Duration
	Object   = zap.Object
	Strings  = zap.Strings
)
