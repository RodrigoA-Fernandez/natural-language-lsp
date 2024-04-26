package rpc_test

import (
	"natural_language_lsp/rpc"
	"testing"
)

type EncodingExample struct {
	Testing bool
}

func TestEncode(t *testing.T) {
	esperado := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	obtenido := rpc.EncodeMessage(EncodingExample{Testing: true})
	if esperado != obtenido {
		t.Fatalf("Expected: %s, Obtenido: %s", esperado, obtenido)
	}
}

func TestDecode(t *testing.T) {
	mensaje := "Content-Length: 15\r\n\r\n{\"Method\":\"hi\"}"
	method, content, err := rpc.DecodeMessage([]byte(mensaje))
  if err != nil {
    t.Fatal(err)
  }
  
  if len(content) != 15 {
    t.Fatalf("Esperado: 15, Devuelto: %d", len(content))
  }

  if method != "hi" {
    t.Fatalf("Esperado: 'hi', Devuelto: %s",method)
  }
}
