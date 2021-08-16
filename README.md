# GitStepCounter

ディレクトリを探索してgitリポジトリを探し、git logコマンドを利用してコードのステップをカウントします。

## Usage

main.goのconfigに条件を入力し、プログラムを実行します。  

```golang:main.go
var authors = []string{"bluemon0919"}
var workspace string = "/Users/kota/go/src/stepcounter/sample"
var since string = "2019-01-01"
var until string = "2021-08-17"
```

```bash
go run *go
```

workspaceで指定したディレクトリをルートとして、サブディレクトリを探索してgitリポジトリを見つけます。  
since-untilで指定した期間にauthorで指定したユーザが変更したコード量を出力します。

標準では次のフォーマットで結果を表示します。
これは全てのユーザー、リポジトリのトータルです。

```bash
: total
: 13 (+10 -3)
```

ユーザ単位、リポジトリ単位で表示したい場合はフラグの設定を変更します。  
ユーザを表示するように設定すると、次のフォーマットで結果を表示します。

```golang
var authorShow = true
```

```bash
: bluemon0919
: 13 (+10 -3)
: sample
: 0 (+0 -0)

: total
: 13 (+10 -3)
```

リポジトリを表示するようにすると、次のフォーマットで結果を表示します。  
リポジトリはユーザごとに表示される仕様です。

```golang
var repositoryShow = false
```

```bash
: bluemon0919
:   /Users/kota/go/src/stepcounter/sample/s1
:   13 (+10 -3)
:   /Users/kota/go/src/stepcounter/sample/s2/s21
:   0 (+0 -0)
:   /Users/kota/go/src/stepcounter/sample/s3
:   0 (+0 -0)
: 13 (+10 -3)
: sample
:   /Users/kota/go/src/stepcounter/sample/s1
:   0 (+0 -0)
:   /Users/kota/go/src/stepcounter/sample/s2/s21
:   0 (+0 -0)
:   /Users/kota/go/src/stepcounter/sample/s3
:   0 (+0 -0)
: 0 (+0 -0)

: total
: 13 (+10 -3)
```
