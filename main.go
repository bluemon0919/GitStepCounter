package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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

// Options represents the operation options of Gitpodem.
type Options struct {
	pull           bool
	authors        Names
	since          string
	until          string
	repositoryShow bool
	authorShow     bool
}

// Gipedom represents an external command being prepared or run.
type Gipedom struct {
	rootDir string
	opt     Options
	w       io.Writer
}

// NewGipedom returns new Gipedom struct.
func NewGipedom(rootDir string, opt Options, w io.Writer) *Gipedom {
	return &Gipedom{
		rootDir: rootDir,
		opt:     opt,
		w:       w,
	}
}

// Run walks the git repository and runs a code step count.
func (g *Gipedom) Run() error {
	paths, err := g.workspaceWalk(g.rootDir)
	if err != nil {
		return err
	}
	if g.opt.pull {
		g.handlePull(paths)
	} else {
		g.handleLog(paths)
	}
	return nil
}

// handlePull handles git pull.
func (g *Gipedom) handlePull(paths []string) {
	ch1 := make(chan []byte, 1)
	for _, dir := range paths {
		go func(dir string) {
			out, _ := g.execPull(dir)
			ch1 <- out
		}(dir)
	}
	for range paths {
		out := <-ch1
		fmt.Fprintln(g.w, string(out))
	}
}

// handleLog handles git log.
func (g *Gipedom) handleLog(paths []string) {
	// -aが指定されていない場合に、全ユーザ対象として検索できるように空データを入れる
	if len(g.opt.authors) == 0 {
		authors = append(authors, "")
	}

	var counts = make(map[string]Count)

	for _, author := range authors {
		if authorShow {
			fmt.Fprintln(g.w, ":", author)
		}
		var authorCount Count

		for _, dir := range paths {
			out, err := g.commandExec(dir, Config{
				name:  author,
				since: since,
				until: until,
			})
			if err != nil {
				return
			}
			c := g.aggregate(out)
			authorCount.add += c.add
			authorCount.sub += c.sub
			authorCount.total += c.total

			if repositoryShow {
				fmt.Fprintln(g.w, ":  ", dir)
				fmt.Fprintf(g.w, ":   %d (+%d -%d)\n", c.total, c.add, c.sub)
			}
		}
		counts[author] = authorCount
		if authorShow {
			fmt.Fprintf(g.w, ": %d (+%d -%d)\n", counts[author].total, counts[author].add, counts[author].sub)
		}
	}

	var t Count
	for _, c := range counts {
		t.add += c.add
		t.sub += c.sub
		t.total += c.total
	}

	fmt.Fprintln(g.w, "")
	fmt.Fprintln(g.w, ": total")
	fmt.Fprintf(g.w, ": %d (+%d -%d)\n", t.total, t.add, t.sub)
}

// commandExec exec git log command.
func (g *Gipedom) commandExec(dir string, c Config) ([]byte, error) {
	args := []string{"-C", dir, "log", "--numstat", "--all", "--author", c.name, "--no-merges", "--pretty=\"%h\""}
	if len(c.since) != 0 {
		args = append(args, "--since", c.since)
	}
	if len(c.until) != 0 {
		args = append(args, "--until", c.until)
	}
	return exec.Command("git", args...).Output()
}

func (g *Gipedom) execPull(dir string) ([]byte, error) {
	args := []string{"-C", dir, "pull"}
	return exec.Command("git", args...).Output()
}

type Count struct {
	add   int
	sub   int
	total int
}

// aggregate aggregates the output data
func (g *Gipedom) aggregate(commnadOutput []byte) Count {

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
			continue
		}
		add, err := strconv.Atoi(d[0])
		if err != nil {
			continue
		}
		sub, err := strconv.Atoi(d[1])
		if err != nil {
			continue
		}
		count.add += add
		count.sub += sub
	}
	count.total = count.add + count.sub
	return count
}

func (g *Gipedom) workspaceWalk(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, f := range files {
		var subdir = filepath.Join(dir, f.Name())
		if f.IsDir() {
			if !g.isGitRepository(subdir) {
				dirs, err := g.workspaceWalk(subdir)
				if err != nil {
					return nil, err
				}
				paths = append(paths, dirs...)
			} else {
				paths = append(paths, subdir)
			}
		}
	}
	return paths, nil
}

func (g *Gipedom) isGitRepository(dir string) bool {
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if f.Name() == ".git" {
			return true
		}
	}
	return false
}
