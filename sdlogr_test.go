package sdlogr_test

import (
	"errors"
	"fmt"
	"github.com/SergeyMosin/sdlogr"
	"github.com/go-logr/logr"
	"testing"
)

func TestPunctuation(t *testing.T) {

	log := sdlogr.New()
	log.Info("") // << not logged
	log.Info("", "k")
	log.Info("", "k", "")
	log.Info("", "", "v")
	log.Info("", "k", "v")
	log.Info("", "k", "v", "k2")
	log.Info("", "k", "v", "k2", "v2")

	fmt.Println("-------------------------------------")

	log.Info("msg")
	log.Info("msg", "k")
	log.Info("msg", "k", "v")
	log.Info("msg", "k", "v", "k2")
	log.Info("msg", "k", "v", "k2", "v2")

	fmt.Println("-------------------------------------")

	log = log.WithName("name1")

	log.Info("") // << not logged
	log.Info("", "k")
	log.Info("", "k", "v")
	log.Info("", "k", "v", "k2")
	log.Info("", "k", "v", "k2", "v2")

	fmt.Println("-------------------------------------")

	log.Info("msg")
	log.Info("msg", "k")
	log.Info("msg", "k", "v")
	log.Info("msg", "k", "v", "k2")
	log.Info("msg", "k", "v", "k2", "v2")

	fmt.Println("-------------------------------------")

	log = log.WithName("name2")

	log.Info("") // << not logged
	log.Info("", "k")
	log.Info("", "k", "v")
	log.Info("", "k", "v", "k2")
	log.Info("", "k", "v", "k2", "v2")

	fmt.Println("-------------------------------------")

	log.Info("msg")
	log.Info("msg", "k")
	log.Info("msg", "k", "v")
	log.Info("msg", "k", "v", "k2")
	log.Info("msg", "k", "v", "k2", "v2")

	fmt.Println("-------------------------------------")

	log = log.WithValues("wvk")

	log.Info("")
	log.Info("", "k", "v")
	log.Info("msg")
	log.Info("msg", "k", "v")

	fmt.Println("-------------------------------------")

	log = log.WithValues("wvk2", "wvv2")

	log.Info("")
	log.Info("", "k", "v")
	log.Info("msg")
	log.Info("msg", "k", "v")

	fmt.Println("-------------------------------------")

	log = log.WithValues("wvk3")

	log.Info("")
	log.Info("", "k", "v")
	log.Info("msg")
	log.Info("msg", "k", "v")

	fmt.Println("-------------------------------------")

	log = sdlogr.New().WithValues("vwk")
	log.Info("")
	log.Info("msg")

	fmt.Println("-------------------------------------")

	log = sdlogr.New().WithValues("vwk", "vwv")
	log.Info("")
	log.Info("msg")
	log = log.WithValues("vwk2")
	log.Info("")
	log.Info("msg")

	fmt.Println("-------------------------------------")

	log = sdlogr.New().WithValues("vwk", "vwv").WithName("name")
	log.Info("")
	log.Info("msg")
	log = log.WithValues("vwk2")
	log.Info("")
	log.Info("msg")
}

func TestKvPairs(t *testing.T) {

	log := sdlogr.New().WithName("kvTester")
	log.Info("int", 42)

	err := errors.New("test_err")
	log.Info("error", err)

	log.Info("nil", nil)

	var b []byte
	log.Info("empty []byte", b)

	var bp *[]byte
	log.Info("*[]byte", bp)

	b = []byte("abc")
	log.Info("[]byte(\"abc\")", b)

	bp = &b
	log.Info("bp=&b", bp)

	var s string
	log.Info("var s string", s)

	var sp *string
	log.Info("var sp *string", sp)

	s = "abc"
	log.Info("s = \"abc\"", s)

	sp = &s
	log.Info("sp = &s", sp)

	st := struct {
		aa string
		pa *string
		bb []byte
		pb *[]byte
	}{
		aa: "aaaa",
	}

	log.Info("struct", st)
	log.Info("&struct", &st)
	log.Info("&struct.bb", &st.bb)
	log.Info("sdlogr.UnmarshalStruct()", sdlogr.UnmarshalStruct(&st))

	m := make(map[string]*string, 2)
	m["aaaaaaaaa"] = &s
	m["bbbbbbbbbbb"] = nil

	log.Info("map", m)
	log.Info("&map", &m)
	log.Info("&map[key]", m["aaaaaaaaa"])

	log.Info("sm", &st, &m)
}

func TestError(t *testing.T) {

	log := sdlogr.New()
	log.Error(nil, "")
	log.Error(nil, "msg")
	log.Error(nil, "msg", "k", "v")

	err := errors.New("test_err")
	log.Error(err, "")
	log.Error(err, "msg")
	log.Error(err, "msg", "k", "v")

	log = log.WithName("name")
	log.Error(err, "")
	log.Error(err, "msg")
	log.Error(err, "msg", "k", "v")

	log = log.WithName("name2").WithValues("wvk", "wvv")
	log.Error(err, "")
	log.Error(err, "msg")
	log.Error(err, "msg", "k", "v")
}

func TestLevel(t *testing.T) {

	log := sdlogr.NewWithOptions(sdlogr.Options{Verbosity: 3})
	log.Info("msg")
	log.V(1).Info("msg V1")
	log.V(2).Info("msg V2")
	log.V(3).Info("msg V3")
	log.V(4).Info("msg V4")
	log.V(5).Info("msg V5")
	log.V(6).Info("msg V6")

	log.V(1).V(1).Info("msq V1.V1")
	log.V(1).V(1).V(1).Info("msq V1.V1.V1")
	log.V(1).V(1).V(1).V(1).Info("msq V1.V1.V1.V1") // << not logged

}

func TestWithValues(t *testing.T) {
	log := sdlogr.New()
	log = log.WithValues("vv", "vvk")
	log.Info("msg1")
	log = log.WithValues("vv", "vvk")
	log.Info("msg2")
	log = log.WithValues("vv1", "vvk")
	log.Info("msg3")
}

func TestOldNew(t *testing.T) {
	log := sdlogr.New()
	log.Info("msg1")

	logNew := log.WithValues("vv", "vvk").WithName("new")
	logNew.Info("logNew msg1")

	err := errors.New("test_err")
	log.Error(err, "msg2")
	logNew.Error(err, "msg2")
}

func TestCallDepth(t *testing.T) {
	log := sdlogr.New().WithCallDepth(0)
	log.Info("cd 0")
	cd1Test(log)
}

func cd1Test(log logr.Logger) {
	log = log.WithValues("f", "cd1Test")
	log.Info("test1")
	log.WithCallDepth(1).Info("test2")
}

func BenchmarkInfo(b *testing.B) {
	log := sdlogr.New().WithName("name").WithValues("wvk", "wvv")
	for n := 0; n < b.N; n++ {
		log.V(5).Info("msg", "k1", "v1", "k2", 10, "k3", nil)
	}
}

func BenchmarkInfoNoCallerInfo(b *testing.B) {
	log := sdlogr.NewWithOptions(sdlogr.Options{LogCallerInfo: false}).
		WithName("name").WithValues("wvk", "wvv")
	for n := 0; n < b.N; n++ {
		log.Info("msg", "k1", "v1", "k2", 10, "k3", nil)
	}
}
