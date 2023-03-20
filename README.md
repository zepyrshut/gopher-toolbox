# excel2struct

Convierte una hoja de excel compatible con la librería [Excelize](https://github.com/qax-os/excelize) a un tipo estructurado de Go. La primera fila debe coincidir con la etiqueta XLSX, sensible a las mayúsculas.

| Id | Nombre | Apellidos | Email                    | Género | Balance |
|----|--------|-----------|--------------------------|--------|---------|
| 1  | Caryl  | Kimbrough | ckimbrough0@fotki.com    | true   | 571.08  |
| 2  | Robin  | Bozward   | rbozward1@thetimes.co.uk | true   | 2162.89 |
| 3  | Tabbie | Kaygill   | tkaygill2@is.gd          | false  | 703.94  |

```go
type User struct {
	Id       int     `xlsx:"Id"`
	Name     string  `xlsx:"Nombre"`
	LastName string  `xlsx:"Apellidos"`
	Email    string  `xlsx:"Email"`
	Gender   bool    `xlsx:"Género"`
	Balance  float32 `xlsx:"Balance"`
}
```

```go
func main() {
 data := exceltostruct.Convert[User]("Book1.xlsx", "Sheet1")
 fmt.Println(data)
}
```

```bash
[{1 Caryl Kimbrough ckimbrough0@fotki.com true 571.08} {2 Robin Bozward rbozward1@thetimes.co.uk true 2162.89} {3 Tabbie Kaygill tkaygill2@is.gd false 703.94}]
```

Donde el primer parámetro es la ruta donde está ubicada la hoja de cálculo y la segunda el nombre de la hoja.

Tipos compatibles: **int**, **float32**, **bool** y **string**.

