# GitStepCounter

ディレクトリを探索してgitリポジトリを探し、git logコマンドを利用してコードのステップをカウントします。  
ルートディレクトリにgitリポジトリを指定することはできません。

## Installation

- 本ツールを利用するためには事前にGitをインストールしてください。
- 本ツールはGoで実装されています。buildに必要なためGoをインストールしてください。

本リポジトリをcloneし、以下のコマンドでインストールします。  
`GOPATH`を設定している場合は`$GOPATH/bin`、設定していない場合は`$HOME/go/bin`がインストール先になります。

```bash
go install
```

## Usage

gipedomコマンドの書式

```bash
gipedom [オプション]
```

次のオプションを用意しています。  
オプション一覧は`-h`をつけて実行することで表示します。

```bash
  -pull
        検索で見つかったgit repositoryに対してpull操作を実行します。
        この操作はカウント機能とは独立して実施します。
  -a value
        GitUserNameを指定します(複数指定可)。
        指定しない場合は全てのユーザが対象となります。
  -d string
        検索のルートディレクトリを指定します。
        指定しない場合は現在のディレクトリが対象となります。
  -s string
        特定の日付より新しいコミットを表示します(2006-01-02)
  -u string
        特定の日付より古いコミットを表示します(2006-01-02)
  -author bool
        ユーザ単位で結果を表示します
  -repos bool
        リポジトリ単位で結果を表示します
```

注意：）複数のAuthorを指定する場合は次のようにします。

```bash
gipedom -a bluemon0919 -a sample
```

ルートディレクトリからサブディレクトリを探索してgitリポジトリを見つけます。  
since, untilで指定した期間にauthorで指定したユーザが変更したコード量を出力します。

標準では次のフォーマットで結果を表示します。
これは全てのユーザー、リポジトリのトータルです。

```bash
: c++
:  0 (+0 -0)
: python
:  16 (+13 -3)
: total
:  16 (+13 -3)
```

ユーザ単位、リポジトリ単位で表示したい場合はオプションを設定します。  
ユーザ単位で表示する場合、次のフォーマットで結果を表示します。

```bash
: author = bluemon0919
: c++
:  0 (+0 -0)
: python
:  16 (+13 -3)
: total
:  16 (+13 -3)

: c++
:  0 (+0 -0)
: python
:  16 (+13 -3)
: total
:  16 (+13 -3)
```

リポジトリ単位で表示する場合、次のフォーマットで結果を表示します。  
リポジトリはユーザごとに表示される仕様です。

```bash
: bluemon0919
:   /home/go/src/stepcounter/sample/s1
:   13 (+10 -3)
:   /home/go/src/stepcounter/sample/s2/s21
:   0 (+0 -0)
:   /home/go/src/stepcounter/sample/s3
:   0 (+0 -0)
: 13 (+10 -3)
: sample
:   /home/go/src/stepcounter/sample/s1
:   0 (+0 -0)
:   /home/go/src/stepcounter/sample/s2/s21
:   0 (+0 -0)
:   /home/go/src/stepcounter/sample/s3
:   0 (+0 -0)
: 0 (+0 -0)

: total
: 13 (+10 -3)
```
