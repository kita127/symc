# symc ![Go](https://github.com/kita127/symc/workflows/Go/badge.svg)

symc is a library that extracts definitions, declarations and references of variables and functions.


## Description

* ライブラリとして機能を提供する
* It serves as a library
* Input preprocessed C source
* モジュール内の以下の情報を解析しデータ化する
* Analyze the following information and convert it to data.
    * Variable definitions
    * Variable declarations
    * Function definitions
    * Function prototype declarations
    * Variables referenced in the function
    * Function call


## Usage

```go
package main

import (
	"fmt"

	"github.com/kita127/symc"
)

func main() {

	cSrc := `
int variable;

int func( void ){

    variable++;

    return 0;
}

`

	module := symc.ParseModule(string(cSrc))
	fmt.Println(module)
}
```

    >go build && ./main
    Module : Statements={ VariableDef : Name=variable, FunctionDef : Name=func, Params=[], Statements=[RefVar : Name=variable] }

## License
This software is released under the MIT License, see LICENSE.
