# Gopher Toolbox

It is a library that gathers the boilerplate code used in different projects. 
It includes the following:

- Database controller implementations:
  - [PGX Pool](github.com/jackc/pgx/v5)
  - MySQL

- Utilities for converting [pgtype](github.com/jackc/pgx/v5/pgtype) to Go types.

- Random data generation for unit tests, similar to the [Faker](https://faker.readthedocs.io/en/master/) 
in Python.

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

- Conversion of Excel files to structured types. You pass the struct type to the
 function `Convert[T any](bookPath, sheetName string)`, and it will return the 
 data as `dataExcel []T`.

- Constants for HTTP handlers.

- Miscelaneous utilities

```go
CorrectTimezone(timeStamp time.Time) time.Time 
GetBool(value string) bool
LogAndReturnError(err error, message string) error
GetBoolFromString(s string) bool
Slugify(s string) string
```