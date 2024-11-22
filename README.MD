# GGenums - Generate Golang enums

## 1. Install

```bash
go install github.com/Gandalf-Le-Dev/ggenums
```

## 2. Usage

You will need to define an enum in a json file like this:

> [!CAUTION]
> The enum file must be named `enums.json` and must be in the same directory as the command.

```json
{
    "package": "my_package",
    "enums": [
        {
            "name": "MyEnum",
            "plural": "MyEnums",
            "values": {
                "value1": "Value1", // You must capitalize the first letter of the value in order to export it, otherwise it will be private
                "value2": "Value2",
                "value3": "Value3"
            }
        },
        {
            "name": "MyOtherEnum",
            "plural": "MyOtherEnums",
            "values": {
                "value1": "Value1",
                "value2": "Value2",
                "value3": "Value3"
            }
        }
    ]
}
```

Then you can generate the enums with the following command:

```bash
ggenums
```

The command must be run in the same directory as the json file.

You can also use golang generate to run the command automatically when you run go generate:

```go
//go:generate ggenums
```

You would typically put the above line in a file called `generate.go` in the same directory as the json file.

## 3. Output

These enums can then be used in your code like this:

```go
package main

import (
    "fmt"
    "my_package"
)

func main() {
    fmt.Println(my_package.MyEnums.Value1) // Output: 0
    fmt.Println(my_package.MyEnums.Value2.String()) // Output: "value2"
    fmt.Println(my_package.MyEnums.Value3) // Output: 2
}
```
