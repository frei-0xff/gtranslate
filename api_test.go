package main

import (
	"strconv"
	"testing"
	"unicode/utf8"

	"github.com/robertkrimen/otto"
)

var vm = otto.New()

func TestCrypt(t *testing.T) {
	type test struct {
		num  string
		op   string
		want string
	}

	tests := []test{
		{num: "461025", op: "+-a^+6", want: "475670226"},
		{num: "475670388", op: "+-a^+6", want: "-2031867735"},
		{num: "-2031867527", op: "+-a^+6", want: "399779148"},
		{num: "399779329", op: "+-a^+6", want: "1775082153"},
		{num: "1775082362", op: "+-a^+6", want: "-1564775049"},
		{num: "-1564774920", op: "+-a^+6", want: "-1842183985"},
		{num: "-1842183776", op: "+-a^+6", want: "1565124006"},
		{num: "1565124136", op: "+-a^+6", want: "-2031023056"},
		{num: "-2031022848", op: "+-a^+6", want: "1242143028"},
		{num: "1242143218", op: "+-a^+6", want: "1907401845"},
		{num: "1907402053", op: "+-a^+6", want: "882159520"},
		{num: "882159698", op: "+-a^+6", want: "-2058833997"},
		{num: "-2058833789", op: "+-a^+6", want: "-1431660943"},
		{num: "-1431660767", op: "+-a^+6", want: "1414727501"},
		{num: "1414727710", op: "+-a^+6", want: "-1561178594"},
		{num: "-1561178451", op: "+-a^+6", want: "1839052455"},
		{num: "1839052487", op: "+-a^+6", want: "-404147828"},
		{num: "-404147619", op: "+-a^+6", want: "-1904186172"},
		{num: "-1904186043", op: "+-a^+6", want: "-1837803712"},
		{num: "-1837803503", op: "+-a^+6", want: "1716177505"},
		{num: "1716177635", op: "+-a^+6", want: "-1825337416"},
	}

	for _, tc := range tests {
		got := crypt(tc.num, tc.op)
		if tc.want != got {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func FuzzCrypt(f *testing.F) {
	testcases := []int32{461025, 475670388, -2031867527, 399779329, 1775082362}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, num int32) {
		res := crypt(strconv.FormatInt(int64(num), 10), "+-a^+6")
		vm.Set("num", num)
		jsRes, _ := vm.Run(`
		function Cn(a, b) {
			for (var c = 0; c < b.length - 2; c += 3) {
				var d = b.charAt(c + 2);
				d = "a" <= d ? d.charCodeAt(0) - 87 : Number(d);
				d = "+" == b.charAt(c + 1) ? a >>> d : a << d;
				a = "+" == b.charAt(c) ? a + d & 4294967295 : a ^ d
			}
			return a
		}
		Cn(num, "+-a^+6")
	`)
		if res != jsRes.String() {
			t.Errorf("Got: %q, JS: %q", res, jsRes)
		}
	})
}

func TestGetTK(t *testing.T) {
	type test struct {
		text string
		ctkk string
		want string
	}

	tests := []test{
		{"000000000", "460914.1766927989", "62714.523400"},
		{"Тестовая строка", "460817.1766927921", "379291.180618"},
		{"Тест Тест Тест", "111111.1111111111", "976525.1004682"},
		{"Т Т Т", "333333.3333333333", "180625.513924"},
	}

	for _, tc := range tests {
		got := getTK(tc.text, tc.ctkk)
		if tc.want != got {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func FuzzGetTk(f *testing.F) {
	f.Add("Тестовая строка", int32(460817), int32(1766927921)) // Use f.Add to provide a seed corpus
	f.Fuzz(func(t *testing.T, text string, t1, t2 int32) {
		if !utf8.ValidString(text) {
			t.Skip()
			return
		}
		for _, v := range []rune(text) {
			if v == 0 || v > 0xFFFF {
				t.Skip()
				return
			}
		}
		ctkk := strconv.FormatInt(int64(t1), 10) + "." + strconv.FormatInt(int64(t2), 10)
		res := getTK(text, ctkk)
		vm.Set("text", text)
		vm.Set("ctkk", ctkk)
		jsRes, _ := vm.Run(`
		function Cn(a, b) {
			for (var c = 0; c < b.length - 2; c += 3) {
				var d = b.charAt(c + 2);
				d = "a" <= d ? d.charCodeAt(0) - 87 : Number(d);
				d = "+" == b.charAt(c + 1) ? a >>> d : a << d;
				a = "+" == b.charAt(c) ? a + d & 4294967295 : a ^ d
			}
			return a
		}
		function Dn(a, b) {
			var c = b.split(".");
			b = Number(c[0]) || 0;
			for (var d = [], e = 0, f = 0; f < a.length; f++) {
				var g = a.charCodeAt(f);
				128 > g ? d[e++] = g : (2048 > g ? d[e++] = g >> 6 | 192 : (55296 == (g & 64512) && f + 1 < a.length && 56320 == (a.charCodeAt(f + 1) & 64512) ? (g = 65536 + ((g & 1023) << 10) + (a.charCodeAt(++f) & 1023),
					d[e++] = g >> 18 | 240, d[e++] = g >> 12 & 63 | 128) : d[e++] = g >> 12 | 224, d[e++] = g >> 6 & 63 | 128), d[e++] = g & 63 | 128)
			}
			a = b;
			for (e = 0; e < d.length; e++) a += d[e], a = Cn(a, "+-a^+6");
			a = Cn(a, "+-3^+b+-f");
			a ^= Number(c[1]) || 0;
			0 > a && (a = (a & 2147483647) + 2147483648);
			c = a % 1E6;
			return c.toString() + "." + (c ^ b)
		}
		Dn(text, ctkk)
	`)
		if res != jsRes.String() {
			t.Errorf("\"%s\" %v %s Got: %q, JS: %q", text, []rune(text), ctkk, res, jsRes)
		}
	})
}
