# symc ![Go](https://github.com/kita127/symc/workflows/Go/badge.svg)
symc is a library that extracts identifiers

symc は 変数名、関数名の定義および参照箇所を抽出するライブラリ


## Description

* ライブラリとして機能を提供する
* プリプロ展開済みのCソースを入力とする
* モジュール内の以下の情報を解析しデータ化する
    * 変数の定義
    * 変数の宣言
    * 関数の定義
    * 関数のプロトタイプ宣言
    * 関数内で参照している変数
    * 関数内での関数コール


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
