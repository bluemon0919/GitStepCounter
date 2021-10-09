package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

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

	var jdata []GitData
	for _, author := range authors {
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

			tmp := g.perse(out)
			for i := range tmp {
				tmp[i].repopath = dir
				tmp[i].author = author
			}
			jdata = append(jdata, tmp...)
		}
	}

	if authorShow {
		adata := make(map[string][]GitData)
		for _, d := range jdata {
			a := adata[d.author]
			a = append(a, d)
			adata[d.author] = a
		}
		for author, d := range adata {
			fmt.Fprintln(g.w, "")
			fmt.Fprintln(g.w, ": author =", author)
			output(g.w, d)
		}
	}
	{
		fmt.Fprintln(g.w, "")
		output(g.w, jdata)
	}
}

func output(w io.Writer, jdata []GitData) {
	kind := make(map[string]AggregateData)
	total := AggregateData{}
	for _, d := range jdata {
		total.add += d.add
		total.sub += d.sub
		k := kind[d.kind]
		k.add += d.add
		k.sub += d.sub
		kind[d.kind] = k
	}
	for key, d := range kind {
		fmt.Fprintln(w, ":", key)
		fmt.Fprintf(w, ":  %d (+%d -%d)\n", d.total(), d.add, d.sub)
	}
	fmt.Fprintln(w, ": total")
	fmt.Fprintf(w, ":  %d (+%d -%d)\n", total.total(), total.add, total.sub)
}

type AggregateData struct {
	add int
	sub int
}

func (d *AggregateData) total() int {
	return d.add + d.sub
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

type GitData struct {
	repopath string
	author   string
	filepath string
	kind     string
	add      int
	sub      int
}

// perse
// コマンドアウトプットを解析してjson形式で返す
// 2   1   test.py
// ↓
// add : 2
// sub : 1
// filepath : test.py
// kind : python
func (g *Gipedom) perse(commnadOutput []byte) (result []GitData) {
	// git logコマンドはoutputにHT(ascii=9)を含み処理しづらいため、SP(ascii=32)に置き換える
	for i := 0; i < len(commnadOutput); i++ {
		c := commnadOutput[i]
		if c == 9 {
			commnadOutput[i] = 32 // HT->SPに置き換え
		}
	}

	slice := strings.Split(string(commnadOutput), string(byte(10))) // LF(ascii=10)で分割する
	for _, str := range slice {
		data := GitData{}
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
		data.add = add
		data.sub = sub
		data.filepath = d[2]
		data.kind = filekind(d[2])
		result = append(result, data)
	}
	return result
}

func filekind(path string) string {
	var kindTable map[string]string = map[string]string{
		".cpp":   "c++",
		".hpp":   "c++",
		".c":     "c",
		".h":     "c",
		".cs":    "c#",
		".css":   "css",
		".dart":  "dart",
		".go":    "go",
		".html":  "html",
		".htm":   "html",
		".java":  "java",
		".js":    "javascript",
		".mat":   "matlab",
		".sql":   "sql",
		".pl":    "perl",
		".php":   "php",
		".py":    "python",
		".rb":    "ruby",
		".rs":    "rust",
		".scala": "scala",
		".sh":    "shellscript",
		".ts":    "typescript",
		".xml":   "xml",
		".vue":   "vuejs",
	}

	e := filepath.Ext(path)
	if val, ok := kindTable[e]; ok {
		return val
	} else {
		return "other"
	}
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
