package main

import (
	// "bufio"
	"encoding/json"
	"fmt"
	"log"
	"natural_language_lsp/analisis"
	"natural_language_lsp/lsp"
	"natural_language_lsp/rpc"
	"natural_language_lsp/scrapper"
	"os"
)

func main() {
	def, err := scrapper.Definir("ornitorrinco")
	if err != nil {
		return
	}

	fmt.Println(def.Palabra)
	fmt.Println(def.Etimologia)
	for _, acepcion := range def.Acepciones {
		fmt.Println(acepcion)
	}

	// logger := getLogger("/home/rodrigo/Proyectos/natural_language_lsp/log.txt")
	// logger.Println("Empecé y eso")
	//
	// scanner := bufio.NewScanner(os.Stdin)
	// scanner.Split(rpc.Split)
	//
	// state := analisis.NewState()
	//
	// for scanner.Scan() {
	// 	msg := scanner.Bytes()
	// 	method, contents, err := rpc.DecodeMessage(msg)
	// 	if err != nil {
	// 		logger.Printf("Error: %s", err)
	// 		continue
	// 	}
	//
	// 	handleMessage(logger, method, contents, state)
	// }
}

func handleMessage(logger *log.Logger, method string, contents []byte, state analisis.State) {
	// logger.Printf("Mensaje recibido con método: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("La petición no se pudo descifrar: %s", err)
			return
		}

		logger.Printf("Connectado a: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)

		// Respuesta
		msg := lsp.NewInitializeResponse(request.ID)
		reply := rpc.EncodeMessage(msg)

		writer := os.Stdout
		writer.Write([]byte(reply))

		logger.Print("Respuesta enviada")
		break

	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("La petición no se pudo descifrar: %s", err)
			return
		}

		state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)

		logger.Printf("textDocument/didOpen: %s",
			request.Params.TextDocument.URI)
		break

	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/didChange: %s", err)
			return
		}

		logger.Printf("Changed: %s",
			request.Params.TextDocument.URI)

		for _, change := range request.Params.ContentChanges {
			state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
		}
		break

	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("Archivo Inválido.")
	}

	return log.New(logfile, "[natural_language_lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
