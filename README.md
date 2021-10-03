# GitStepCounter

ディレクトリを探索してgitリポジトリを探し、git logコマンドを利用してコードのステップをカウントします。  
ルートディレクトリにgitリポジトリを指定することはできません。

## Usage

main.goのconfigに条件を入力し、プログラムを実行します。  

```bash
go run *go
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
go run *go -a bluemon0919 -a sample
```

ルートディレクトリからサブディレクトリを探索してgitリポジトリを見つけます。  
since, untilで指定した期間にauthorで指定したユーザが変更したコード量を出力します。

標準では次のフォーマットで結果を表示します。
これは全てのユーザー、リポジトリのトータルです。

```bash
: total
: 13 (+10 -3)
```

ユーザ単位、リポジトリ単位で表示したい場合はオプションを設定します。  
ユーザ単位で表示する場合、次のフォーマットで結果を表示します。

```bash
: bluemon0919
: 13 (+10 -3)
: sample
: 0 (+0 -0)

: total
: 13 (+10 -3)
```

リポジトリ単位で表示する場合、次のフォーマットで結果を表示します。  
リポジトリはユーザごとに表示される仕様です。

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
