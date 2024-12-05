package day19

import (
	"adventofcode/m/v2/util"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func Day19(inputFile string, part int) {
	if part == 0 {
		fmt.Printf("Test: %d\n", RunAll(inputFile))
	} else {
		fmt.Println("Not implmenented.")
	}
}

func RunAll(inputFile string) int {
	ls := util.LineScanner(inputFile)
	line, ok := util.Read(ls)

	// The Program c is a set of workflows
	c := NewProgram()
	for line != "" {
		c.addWorkflow(line)
		line, ok = util.Read(ls)
	}

	// Read parts
	line, ok = util.Read(ls)
	for ok {
		c.addPart(line)
		line, ok = util.Read(ls)
	}

	sum := 0
	for _, p := range c.parts {
		if c.workflows["in"].run(p) {
			sum += p.sum()	
		}
	}

	return sum 
}

// =======
// PROGRAM
// =======
type Program struct {
	workflows map[string]*Workflow
	parts []*Part
}

func NewProgram() *Program {
	workflows := make(map[string]*Workflow)
	workflows["A"] = &Workflow{"A", nil, "Accept"}
	workflows["R"] = &Workflow{"R", nil, "Reject"}
	
	parts := []*Part{}
 return &Program{workflows, parts}
}

func (c *Program) addWorkflow(line string) *Workflow {
	//  The workflow is either a known, named workflow, where line == workflow id
	if strings.Index(line, ":") == -1 {
		return c.workflows[line]
	}

	var l func(p *Part) *Workflow
	var name, d string
	// ... or a new, named workflow, where line == (id){ [x|m|a|s] <|> n: w1, w2}
	named := regexp.MustCompile(`([a-z]+)\{(x|m|a|s)(<|>)(\d+):(.*?),(.*)}`)
	m := named.FindAllStringSubmatch(line, 10)
	if len(m) == 0 {
		// ... or a new, unnamed workflow, where line == [x|m|a|s] <|> n: w1, w2
		unnamed := regexp.MustCompile(`(x|m|a|s)(<|>)(\d+):(.*?),(.*)`)
		m = unnamed.FindAllStringSubmatch(line, 10)
		l, d = c.NewLambda(m[0][1], m[0][2], m[0][3] ,m[0][4], m[0][5])
	}	else {
		name = m[0][1]
		l, d = c.NewLambda(m[0][2], m[0][3], m[0][4], m[0][5], m[0][6])
	}

	w := &Workflow{
		id: name,
		next: l, 
		descr: d,
	}
	if w.id != "" {
		c.workflows[w.id] = w
	}

	return w
}

func (c *Program) addPart(line string) {
	reg := regexp.MustCompile(`\{x=(\d+),m=(\d+),a=(\d+),s=(\d+)\}`)
	m := reg.FindAllStringSubmatch(line, 5)
	X, _ := strconv.Atoi(m[0][1])
	M, _ := strconv.Atoi(m[0][2])
	A, _ := strconv.Atoi(m[0][3])
	S, _ := strconv.Atoi(m[0][4])
	c.parts = append(c.parts, NewPart(X, M, A, S))
}

// =====
// PARTS
// =====
type Part struct {
	p map[string]int
}

func NewPart(x, m, a, s int) *Part {
	return &Part{p: map[string]int{"x": x, "m": m, "a": a, "s": s}}
}

func (p *Part) get(s string) int {
	return p.p[s]
}

func (p *Part) sum() int {
	return p.p["x"]+p.p["m"]+p.p["a"]+p.p["s"]
}

// =========
// WORKFLOWS
// =========
type Workflow struct {
	id 	string
	next func(p *Part) *Workflow
	descr string
}

func (w *Workflow) terminated() bool {
	if w.id == "A" || w.id == "R" {
		return true
	}
	return false
}

func (w *Workflow) run(p *Part) bool {
	for !w.terminated() {
		w = w.next(p) 
	}	
	
	return w.id == "A"
}

func (c *Program) NewLambda(letter string, op string, num string, then string, orElse string) (func(p *Part) *Workflow, string) {
	return func(p *Part) *Workflow {
			if op == "<" {
				if p.get(letter) < asInt(num) {
					return c.addWorkflow(then)
				} else {
					return c.addWorkflow(orElse)
				}
			} else {
				if p.get(letter) > asInt(num) {
					return c.addWorkflow(then)
				} else {
					return c.addWorkflow(orElse)
				}
			}
		},
		fmt.Sprintf("%s%s%s:%s:%s", letter, op, num, then, orElse)

	}

// UTIL
func asInt(str string) int {
	x, _ := strconv.Atoi(str)
	return x
}