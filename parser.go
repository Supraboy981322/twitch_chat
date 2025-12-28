package main

import (
	"strings"
	"strconv"
)

type (
	ParserMem struct {
		Do bool
		Typ string
		Mem string
		Reset func()
		Cancel func()
	}
	StrPar struct {
		Pos int
		Cur string
		Out string
		Len int
		In []string
		Defs map[string]string
		Parser ParserMem
	}
)
var (
)

func parser(in string) string {
	pMem := ParserMem {
		Do: false,
		Typ: "",
		Mem: "",
	}
	p := StrPar{
		Pos: 0,
		Cur: "",
		Out: "",
		Len: len(in),
		In: strings.Split(in, ""),
		Defs: defs,
		Parser: pMem,
	}
	p.Parser.Reset = func(){
			p.Parser.Mem = ""
			p.Parser.Do = false
			p.Parser.Typ = ""
	};p.Parser.Cancel = func(){
		p.Parser.Reset()
		p.Len = 0
	}
	return p.next()
}

func (p *StrPar) next() string {
	if p.eof() { return p.Out }
	
	p.Cur = p.In[p.Pos]

	if p.Parser.Do {
		p.eval()
		p.Pos++
		if p.eof() { return p.Out }
		return p.next()
	}

	switch p.Cur {
	 case "6":
		p.Parser.Typ = "cringe;67"
		p.Parser.Do = true
	 case "@":
		p.Parser.Do = true
		p.Parser.Typ = "at"
	 case ".":
		p.Out += p.Cur
		p.checkLink()
   default:
		p.Out += p.Cur
	}
	p.Pos++

	return p.next()
}

func (p *StrPar) eof() bool {
	return p.Pos >= p.Len
}

func (p *StrPar) eval() {
	if p.eof() {
		switch p.Parser.Typ {
		 case "at":
			print(p.Cur)
			p.Out += "\033[1;34m"+p.Parser.Mem+"\033[0m"
			p.eat()
			p.Parser.Reset()
		 default:
		}
		return
	}
	switch p.Cur {
	 case "7":
		if p.Parser.Typ == "cringe;67" {
			p.Out += "[cringe removed]"
			p.eat()
			p.Parser.Reset()
		}
   case " ":
		if p.Parser.Typ == "at" {
			p.Out += "\033[1;30;44m@"+p.Parser.Mem+"\033[0m "
			p.Parser.Reset()
		}
	 default:
		if p.Parser.Typ == "cringe;67" {
			switch p.Cur {
			 case ":", "|", ",", " ", "_", "-", "6": fallthrough
			 case ".", "/", "\\", "!", "\t", "*", "&", "=", "+": p.eat()
			 default:
				p.Out += p.Parser.Mem
				p.Parser.Reset()
			}
		} else {
			p.Parser.Mem += p.Cur
		}
	}
}

func (p *StrPar) eat() {
	if p.eof() { return }
	p.Pos++ 
}

func (p *StrPar) checkLink() {
	restR := strings.Join(p.In[p.Pos+1:], "")
	restS := strings.FieldsFunc(restR, func(r rune) bool {
		switch r {
		 case ' ', '\n': return true
		 default: return false
		}
	})
	var ext string
	if len(restS) > 0 {
		ext = restS[0]
	} else { return }
	
	if ext == "" { return }
	
	if isLink, onLine := chkTLD(ext); isLink {
		p.Out = "\033[3;5;31m[CHAT NOT DISPLAYED "+
				"(detected TLD from line "+strconv.Itoa(onLine)+" of "+
				"the offical IANA tld list)]\033[0m"
		p.Parser.Cancel()
	}
}
