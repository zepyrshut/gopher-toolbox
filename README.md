# Gopher Toolbox

Es una librería donde se concentra el código _boilerplate_ que se usan en
distintos proyectos. Incluyen lo siguiente:

- Implementación controladores de bases de datos:
  - [PGX Pool](github.com/jackc/pgx/v5)
  - MySQL


- Utilidades para conversión de tipos [pgtype](github.com/jackc/pgx/v5/pgtype) a
tipos de Golang.


- Generación de datos aleatorios para pruebas unitarias, similar a librería 
[Faker](https://faker.readthedocs.io/en/master/) de Python.

```go
MaleName() string
FemaleName() string
Name() string
LastName() string
Email(beforeAt string) string
Int(min, max int64) int64
Float(min, max float64) float64
Bool() bool
Chars(min, max int) string
AllChars(min, max int) string
AllCharsOrEmpty(min, max int) string
AllCharsOrNil(min, max int) *string
NumericString(length int) string
Sentence(min, max int) string
```

- Conversión de ficheros Excel a tipos estructurados. Se le pasa el tipo del 
_struct_ a la función `Convert[T any](bookPath, sheetName string)` y te
devolverá los datos del tipo `dataExcel []T`.


- Constantes para los manejadores HTTP.


- Utilidades varias

```go
CorrectTimezone(timeStamp time.Time) time.Time 
GetBool(value string) bool
LogAndReturnError(err error, message string) error
GetBoolFromString(s string) bool
Slugify(s string) string
```