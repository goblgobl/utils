package log

import "io"

type Noop struct {
}

func (_ Noop) Log()                                   {}
func (_ Noop) LogTo(io.Writer)                        {}
func (_ Noop) Reset()                                 {}
func (_ Noop) Release()                               {}
func (_ Noop) Bytes() []byte                          { return nil }
func (n Noop) Info(ctx string) Logger                 { return n }
func (n Noop) Warn(ctx string) Logger                 { return n }
func (n Noop) Error(ctx string) Logger                { return n }
func (n Noop) Fatal(ctx string) Logger                { return n }
func (n Noop) Request(route string) Logger            { return n }
func (n Noop) Err(err error) Logger                   { return n }
func (n Noop) Int(key string, value int) Logger       { return n }
func (n Noop) Int64(key string, value int64) Logger   { return n }
func (n Noop) String(key string, value string) Logger { return n }
func (n Noop) Binary(key string, value []byte) Logger { return n }
func (n Noop) Bool(key string, value bool) Logger     { return n }
func (n Noop) Field(field Field) Logger               { return n }
func (n Noop) Fixed()                                 { return }
func (n Noop) MultiUse() Logger                       { return n }
