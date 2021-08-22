package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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

var authors Names    // Git author names {"bluemon0919", "sample"}
var workspace string // Root directory= "/Users/kota/go/src/stepcounter/sample"
var since string     // "2006-01-02"
var until string     // "2006-01-02"

var repositoryShow = true
var authorShow = true

func init() {
	flag.Var(&authors, "a", "Authorを指定します(複数指定可)。指定しない場合は全てのユーザが対象となります。")
	flag.StringVar(&workspace, "d", "", "検索のルートディレクトリを指定します。指定しない場合は現在のディレクトリが対象となります。")
	flag.StringVar(&since, "s", "", "特定の日付より新しいコミットを表示します(2006-01-02)")
	flag.StringVar(&until, "u", "", "特定の日付より古いコミットを表示します(2006-01-02)")
}
func main() {
	flag.Parse()
	workspace, err := filepath.Abs(workspace)
	if err != nil {
		panic(err)
	}
	fmt.Println("authors:", authors)
	fmt.Println("rootDir:", workspace)
	fmt.Println("since:", since)
	fmt.Println("until:", until)

	paths := workspaceWalk(workspace)
	fmt.Println(paths)

	// -aが指定されていない場合に、全ユーザ対象として検索できるように空データを入れる
	if len(authors) == 0 {
		authors = append(authors, "")
	}

	var counts = make(map[string]Count)

	for _, author := range authors {
		if authorShow {
			fmt.Println(":", author)
		}
		var authorCount Count

		for _, dir := range paths {
			out := commandExec(dir, Config{
				name:  author,
				since: since,
				until: until,
			})
			//fmt.Println(string(out))
			c := aggregate(out)
			authorCount.add += c.add
			authorCount.sub += c.sub
			authorCount.total += c.total

			if repositoryShow {
				fmt.Println(":  ", dir)
				fmt.Printf(":   %d (+%d -%d)\n", c.total, c.add, c.sub)
			}
		}
		counts[author] = authorCount
		if authorShow {
			fmt.Printf(": %d (+%d -%d)\n", counts[author].total, counts[author].add, counts[author].sub)
		}
	}

	var t Count
	for _, c := range counts {
		t.add += c.add
		t.sub += c.sub
		t.total += c.total
	}

	fmt.Println("")
	fmt.Println(": total")
	fmt.Printf(": %d (+%d -%d)\n", t.total, t.add, t.sub)
}

// commandExec exec git log command.
func commandExec(dir string, c Config) []byte {
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	var cmd *exec.Cmd
	args := []string{"log", "--numstat", "--all", "--author", c.name, "--no-merges", "--pretty=\"%h\""}
	if len(c.since) != 0 {
		args = append(args, "--since", c.since)
	}
	if len(c.until) != 0 {
		args = append(args, "--until", c.until)
	}
	cmd = exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return out
}

type Count struct {
	add   int
	sub   int
	total int
}

// aggregate aggregates the output data
func aggregate(commnadOutput []byte) Count {

	// git logコマンドはoutputにHT(ascii=9)を含み処理しづらいため、SP(ascii=32)に置き換える
	for i := 0; i < len(commnadOutput); i++ {
		c := commnadOutput[i]
		if c == 9 {
			commnadOutput[i] = 32 // HT->SPに置き換え
		}
	}

	var count Count
	slice := strings.Split(string(commnadOutput), string(byte(10))) // LF(ascii=10)で分割する
	for _, str := range slice {
		if len(str) == 0 { // 空の文字列は除外
			continue
		}
		if []byte(str)[0] == 34 { // ダブルクォーテーションから始まる文字列は除外
			continue
		}

		d := strings.Split(str, string(" "))
		if len(d) != 3 {
			log.Fatal(str)
			continue
		}
		add, err := strconv.Atoi(d[0])
		if err != nil {
			panic(err)
		}
		sub, err := strconv.Atoi(d[1])
		if err != nil {
			panic(err)
		}
		count.add += add
		count.sub += sub
		//fmt.Printf("[%s]\n", str)
	}
	count.total = count.add + count.sub
	return count
}

func workspaceWalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, f := range files {
		var subdir = filepath.Join(dir, f.Name())
		if f.IsDir() {
			if !isGitRepository(subdir) {
				paths = append(paths, workspaceWalk(subdir)...)
			} else {
				paths = append(paths, subdir)
				//fmt.Println(subdir)
			}
		}
	}
	return paths
}

func isGitRepository(dir string) bool {
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if f.Name() == ".git" {
			return true
		}
	}
	return false
}
