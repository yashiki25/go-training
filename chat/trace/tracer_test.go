package trace

import (
	"bytes"
	"testing"
)

// tips: 名前がTestではじまり、*testing.T型の引数を受け取るメソッドはユニットテストとみなされる
func TestNew(t *testing.T)  {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	} else {
		tracer.Trace("こんにちは、traceパッケージです")
		if buf.String() != "こんにちは、traceパッケージです\n" {
			t.Errorf("'%s'という誤った文字列が出力されました", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("データ")
}