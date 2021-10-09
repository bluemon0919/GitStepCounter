package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// git logコマンド実行時の設定
type Config struct {
	name  string // author
	since string
	until string
}

type Names []string

func (ns *Names) String() string {
	return fmt.Sprintf("%s", *ns)
}

func (ns *Names) Set(value string) error {
	*ns = append(*ns, value)
	return nil
}

var pull bool        // git pull exec flag
var authors Names    // Git author names {"bluemon0919", "sample"}
var workspace string // Root directory= "/Users/kota/go/src/stepcounter/sample"
var since string     // "2006-01-02"
var until string     // "2006-01-02"

var repositoryShow bool // Include repository information in result.
var authorShow bool     // Include author in result.

func init() {
	flag.BoolVar(&pull, "pull", false, "ルートディレクトリ以下のリポジトリに対してgit pull操作を行います。")
	flag.Var(&authors, "a", "Authorを指定します(複数指定可)。指定しない場合は全てのユーザが対象となります。")
	flag.StringVar(&workspace, "d", "", "検索のルートディレクトリを指定します。指定しない場合は現在のディレクトリが対象となります。")
	flag.StringVar(&since, "s", "", "特定の日付より新しいコミットを表示します(2006-01-02)")
	flag.StringVar(&until, "u", "", "特定の日付より古いコミットを表示します(2006-01-02)")
	flag.BoolVar(&repositoryShow, "repos", false, "リポジトリ単位で結果を表示します")
	flag.BoolVar(&authorShow, "author", false, "ユーザ単位で結果を表示します")
}

func main() {
	flag.Parse()
	workspace, err := filepath.Abs(workspace)
	if err != nil {
		panic(err)
	}
	fmt.Println("pull:", pull)
	fmt.Println("authors:", authors)
	fmt.Println("rootDir:", workspace)
	fmt.Println("since:", since)
	fmt.Println("until:", until)
	fmt.Println("author:", authorShow)
	fmt.Println("repos:", repositoryShow)

	pedom := NewGipedom(workspace, Options{
		pull:           pull,
		authors:        authors,
		since:          since,
		until:          until,
		authorShow:     authorShow,
		repositoryShow: repositoryShow,
	}, os.Stdout)
	if err := pedom.Run(); err != nil {
		fmt.Print(err)
	}
}
