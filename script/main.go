package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type FileChanges struct {
	Current string `json:"current"`
	Changes string `json:"changes"`
}

type RequestData struct {
	Files map[string]FileChanges `json:"files"`
}

func main() {
	// Inicializa el mapa de archivos
	fileChanges := make(map[string]FileChanges)

	// Obtiene la lista de archivos cambiados desde el git
	output, err := exec.Command("git", "diff", "--name-only", "HEAD^", "HEAD").Output()
	if err != nil {
		fmt.Printf("Error al obtener archivos cambiados: %v\n", err)
		return
	}

	log.Println("obtener", output)

	files := strings.Split(string(output), "\n")
	for _, file := range files {
		if file == "" {
			continue
		}

		// Obtiene el contenido actual y los cambios del archivo
		currentContent, err := exec.Command("git", "show", "HEAD:"+file).Output()
		if err != nil {
			fmt.Printf("Error al obtener contenido actual de %s: %v\n", file, err)
			continue
		}

		log.Println("current", currentContent)

		changes, err := exec.Command("git", "diff", "--unified=0", "HEAD^", "HEAD", "--", file).Output()
		if err != nil {
			fmt.Printf("Error al obtener cambios de %s: %v\n", file, err)
			continue
		}

		log.Println("changes", changes)

		if len(changes) == 0 {
			fmt.Printf("No changes detected for file: %s\n", file)
		} else {
			fmt.Printf("obtener changes: %s\n", string(changes))
		}

		fileChanges[file] = FileChanges{
			Current: string(currentContent),
			Changes: string(changes),
		}
	}

	jsonData, err := json.Marshal(fileChanges)
	if err != nil {
		fmt.Printf("Error al convertir a JSON: %v\n", err)
		return
	}

	log.Println("json", string(jsonData))

	fmt.Printf("Respuesta de la LLM: %s\n", "hi")
}
