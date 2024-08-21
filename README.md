# gorender

Simple y minimalista librería para procesar plantillas utilizando la librería
estándar de Go `html/template`.

## Características

- Procesamiento de plantillas utilizando `html/template`.
- Soporte para caché de plantillas.
- Posibilidad de añadir funciones personalizadas a las plantillas.
- Configuración sencilla con opciones por defecto que se pueden sobreescribir.
Inspirado en `Gin`.

## Instalación

```bash
go get github.com/zepyrshut/gorender
```

## Uso mínimo

Las plantillas deben tener la siguiente estructura, observa que las páginas a
procesar están dentro de `pages`. Los demás componentes como bases y fragmentos
pueden estar en el directorio raíz o dentro de un directorio.

Puedes cambiar el nombre del directorio `template` y `pages`. Ejemplo en la
siguiente sección.

```
template/
├── pages/
│   └── page.html
├── base.html
└── fragment.html
```

```go
import (
    "github.com/zepyrshut/gorender"
)

func main() {
    ren := gorender.New()

    // ...

    td := &gorender.TemplateData{}
    ren.Template(w, r, "index.html", td)

    // ...
}
```

## Personalización

> Recuerda que si habilitas el caché, no podrás ver los cambios que realices
> durante el desarrollo.

```go
func dummyFunc() string {
	return "dummy"
}

func main() {

    customFuncs := template.FuncMap{
		"dummyFunc": dummyFunc,
	}

    renderOpts := &gorender.Render{
		EnableCache:       true,
		TemplatesPath:     "template/path",
		PageTemplatesPath: "template/path/pages",
		Functions:         customFuncs,
	}

    ren := gorender.New(gorender.WithRenderOptions(renderOpts))

    // ...

    td := &gorender.TemplateData{}
    ren.Template(w, r, "index.html", td)

    // ...
}
```
## Agradecimientos

- [Protección CSRF justinas/nosurf](https://github.com/justinas/nosurf)
- [Valicación go-playground/validator](https://github.com/go-playground/validator)

## Descargo de responsabilidad

Esta librería fue creada para usar las plantillas en mis proyectos privados, es 
posible que también solucione su problema. Sin embargo, no ofrezco ninguna 
garantía de que funcione para todos los casos de uso, tenga el máximo
rendimiento o esté libre de errores.

Si decides integrarla en tu proyecto, te recomiendo que la pruebes para 
asegurarte de que cumple con tus expectativas y requisitos.

Si encuentras problemas o tienes sugerencias de mejora, puedes colocar tus 
aportaciones a través de _issues_ o _pull requests_ en el repositorio. Estaré 
encantado de ayudarte.


