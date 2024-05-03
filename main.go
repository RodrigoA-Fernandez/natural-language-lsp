package main

import (
	"bufio"
	"encoding/json"
	"io"

	// "fmt"
	"log"
	"natural_language_lsp/analisis"
	"natural_language_lsp/lsp"
	"natural_language_lsp/rpc"
	"os"
)

func main() {
	logger := getLogger("/home/rodrigo/Proyectos/natural_language_lsp/log.txt")
	logger.Println("Empecé y eso")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analisis.NewState()
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Error: %s", err)
			continue
		}

		handleMessage(logger, writer, method, contents, state)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, method string, contents []byte, state analisis.State) {
	logger.Printf("Mensaje recibido con método: %s", method)

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
		writeResponse(writer, lsp.NewInitializeResponse(request.ID))

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

	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/hover: %s", err)
			return
		}

		logger.Printf(string(contents))

		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position, logger)

		writeResponse(writer, response)

	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("Archivo Inválido.")
	}

	return log.New(logfile, "[natural_language_lsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
