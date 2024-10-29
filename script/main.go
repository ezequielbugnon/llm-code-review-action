package main

import (
	"fmt"
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

		changes, err := exec.Command("git", "diff", "HEAD^", "HEAD", "--", file).Output()
		if err != nil {
			fmt.Printf("Error al obtener cambios de %s: %v\n", file, err)
			continue
		}

		fileChanges[file] = FileChanges{
			Current: string(currentContent),
			Changes: string(changes),
		}
	}

	// Prepara los datos para enviar a la LLM
	/*requestData := RequestData{Files: fileChanges}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Printf("Error al convertir a JSON: %v\n", err)
		return
	}*/

	// Envía los datos a la API de la LLM
	/*resp, err := http.Post("URL_DE_LA_API_DE_LLM", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error al enviar la solicitud: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error al leer la respuesta: %v\n", err)
		return
	}

	fmt.Printf("Respuesta de la LLM: %s\n", body)*/

	fmt.Printf("Respuesta de la LLM: %s\n", fileChanges)
}
