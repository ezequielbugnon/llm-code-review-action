# LLM Comparison

Este proyecto utiliza un script en Go para realizar análisis de código y generar comentarios en pull requests en GitHub. El análisis se realiza a través de una API de revisión de código.

## Dependencias

- Go 1.23 o superior

## Uso
Este proyecto está diseñado para ejecutarse automáticamente mediante GitHub Actions. Se activa cuando se abre o sincroniza un pull request. El flujo de trabajo realiza lo siguiente:

- Checkout del código: Obtiene el código del repositorio.
- Configuración de Go: Configura el entorno de Go.
- Instalación de dependencias: Ejecuta go mod tidy para gestionar las dependencias.
- Ejecución del analizador: Ejecuta el script principal (main.go) que realiza la - revisión del código utilizando la API de StackSpot.
- Publicación de la revisión: Publica el resultado del análisis como un comentario en el pull request.

## Variables de Entorno
Asegúrate de configurar las siguientes variables de entorno en los secretos de tu repositorio de GitHub:

-CLIENTID: Tu Client ID para la API de StackSpot.
-CLIENTSECRET: Tu Client Secret para la API de StackSpot.
-GITHUB_TOKEN: Token de acceso para publicar comentarios en el pull request.

## Contribuciones
Las contribuciones son bienvenidas. Si deseas colaborar, por favor abre un pull request o un issue para discutir tus ideas.

