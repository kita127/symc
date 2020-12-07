# symc

symc is a library that extracts identifiers

### Input

プリプロ展開済みのCソース

### Description

* ライブラリとして機能を提供する
* プリプロ展開済みのCソースを入力とする
* モジュール内の関数が参照する識別子の情報をデータ化して提供する
* データの内容は以下
    * モジュールが定義する変数情報
        * 識別子名
    * モジュールが定義する関数情報
        * 関数名
        * 関数が参照する変数
        * 関数が参照する関数


### Usage

```go

src := `
extern ex_var;
int main(void){
    ex_var++;
}
`

data := ModName.Parse(src)

```
