function dec2bin(dec) {
  return (dec >>> 0).toString(2);
}
function Cn(a, b) {
  for (var c = 0; c < b.length - 2; c += 3) {
    var d = b.charAt(c + 2);
    d = "a" <= d ? d.charCodeAt(0) - 87 : Number(d);
    d = "+" == b.charAt(c + 1) ? a >>> d : a << d;
    a = "+" == b.charAt(c) ? a + d & 4294967295 : a ^ d
  }
  return a
}
function Dn(text, ctkk) {
  console.log(text, ctkk)
  var parts = ctkk.split(".");
  var t = Number(parts[0]) || 0;
  for (var buf = [], j = 0, i = 0; i < text.length; i++) {
    var ch = text.charCodeAt(i);
    128 > ch ? buf[j++] = ch : (2048 > ch ? buf[j++] = ch >> 6 | 192 : (55296 == (ch & 64512) && i + 1 < text.length && 56320 == (text.charCodeAt(i + 1) & 64512) ? (ch = 65536 + ((ch & 1023) << 10) + (text.charCodeAt(++i) & 1023),
      buf[j++] = ch >> 18 | 240, buf[j++] = ch >> 12 & 63 | 128) : buf[j++] = ch >> 12 | 224, buf[j++] = ch >> 6 & 63 | 128), buf[j++] = ch & 63 | 128)
  }
  text = t;
  for (j = 0; j < buf.length; j++){
    text += buf[j];
    text = Cn(text, "+-a^+6");
  }
  text = Cn(text, "+-3^+b+-f");
  console.log(text)
  text ^= Number(parts[1]) || 0;
  console.log(text,dec2bin(text),dec2bin(text & 2147483647),dec2bin((text & 2147483647) + 2147483648))
  0 > text && (text = (text & 2147483647) + 2147483648);
  parts = text % 1E6;
  return parts.toString() + "." + (parts ^ t)
}

// console.log(Dn("Тестовая строка", "460817.1766927921"))
// console.log(Dn("Тест Тест Тест", "111111.1111111111"))
// console.log(Dn("Т Т Т", "333333.3333333333"))
console.log(Dn("000000000", "460914.1766927989"))
// Cn(461025,"+-a^+6")
