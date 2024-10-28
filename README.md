# Golang web app starter template

This project is made to accelerate a web app creation by handling:
- standard project structure
- server creation with secured parameters
- errors with middlware
- personalized errors (and even forcing error for testing purpose)
- html templates layout
- dark and light theme with system preference detection and toggle button
- responsive
- environment variables (for secrets)
- hooks with NTFY.sh (for sending server alerts) (see .env)
- and more (css, favicon, etc.)

## Usage

```sh
# rename and fill the .env_example to .env 
# clone the repo
git clone https://github.com/nicgen/golang-web-starter.git <your project>
cd <your project>
```

run it

```go
// with air
air
// manualy
go run cmd/groupie_tracker/main.go
```

Note: [link](https://hello-there.org/go/tools/live-server/) for installing **air**

## Architecture

```txt
├── cmd
├── .env
├── go.mod
├── handlers
├── models
├── README.md
├── static
│   ├── css
│   ├── img
│   └── js
└── templates
```

### Error handling

> Middleware refers to a function that wraps an HTTP handler to add additional behavior before or after the handler processes an HTTP request  

The middleware uses HandleError within a defer statement to handle any panics that occur during the request processing. It logs the panic and sends an appropriate error response using HandleError.  

The HandleError function will set the appropriate HTTP status code and render an error page with the provided message.  
call HandleError directly within your HTTP handlers when you encounter an error.  

### Single go.mod

**Simplicity**: A single go.mod file at the root of your project keeps things simple and straightforward. It reduces complexity by centralizing dependency management.

**Consistency**: Ensures that all parts of the project use the same versions of dependencies, avoiding conflicts and inconsistencies.

**Standard Practice**: This follows the convention used in most Go projects, making it easier for other developers to understand and navigate your project.

<!-- todo, add mux -->

<!--
### Reproducibility

```bash
mkdir -p {cmd,handlers,models,static/{css,img,js},templates}
mkdir cmd/<project-name>
touch cmd/<project-name>/main.go handlers/{index,about,error}.go README.md static/{css/styles.css,img/about.txt} templates/{about,error,index,layout}.html
go mod init <project-name>
air init
sed 's/  cmd = "go build -o .\/tmp\/main ."/  cmd = "go build -o .\/tmp\/main .\/cmd\/<project-name>\/main.go"/g' .air.toml
air -c .air.toml
# after that launch it with `air`
```
-->

## Attribution

- The icons for the toggle buton are from [Tabler Icons](https://tabler.io/icons) (MIT)  
- The SVG icons were optimized with [SVGOMG](https://jakearchibald.github.io/svgomg/) powered by [SVG Optimizer](https://github.com/svg/svgo) (MIT)  
- The SVG were encoded by [Url encoder for SVG](https://yoksel.github.io/url-encoder/)  

The favicon was generated on [favicon.io](https://favicon.io/) using the following graphics from Twitter Twemoji:  
- Graphics Title: 2620.svg  
- Graphics Author: Copyright 2020 Twitter, Inc and other contributors (https://github.com/twitter/twemoji)  
- Graphics Source: https://github.com/twitter/twemoji/blob/master/assets/svg/2620.svg  
- Graphics License: CC-BY 4.0 (https://creativecommons.org/licenses/by/4.0/)  
