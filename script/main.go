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
		log.Println("Error al obtener archivos cambiados: ", err)
		return
	}

	log.Println("obtener", string(output))

	files := strings.Split(string(output), "\n")
	for _, file := range files {
		if file == "" {
			continue
		}

		// Obtiene el contenido actual y los cambios del archivo
		currentContent, err := exec.Command("git", "show", "HEAD:"+file).Output()
		if err != nil {
			log.Println("Error al obtener contenido actual de ", file, err)
			continue
		}

		log.Println("current", string(currentContent))

		changes, err := exec.Command("git", "diff", "--unified=0", "HEAD^", "HEAD", "--", file).Output()
		if err != nil {
			log.Println("Error al obtener cambios de ", file, err)
			continue
		}

		log.Println("changes", string(changes))

		fileChanges[file] = FileChanges{
			Current: string(currentContent),
			Changes: string(changes),
		}
	}

	jsonData, err := json.Marshal(fileChanges)
	if err != nil {
		log.Println("Error al convertir a JSON:", err)
		return
	}

	log.Println("json", string(jsonData))

	fmt.Println("hi")
}
