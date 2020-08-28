package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os/user"
	"path"
	"time"
)

const (
	BookOfMormonURL         = "https://github.com/bcbooks/scriptures-json/raw/master/book-of-mormon.json"
	DoctrineAndCovenantsURL = "https://github.com/bcbooks/scriptures-json/raw/master/doctrine-and-covenants.json"
	NewTestamentURL         = "https://github.com/bcbooks/scriptures-json/raw/master/new-testament.json"
	OldTestamentURL         = "https://github.com/bcbooks/scriptures-json/raw/master/old-testament.json"
	PearlOfGreatPriceURL    = "https://github.com/bcbooks/scriptures-json/raw/master/pearl-of-great-price.json"
)

type Scriptures struct {
	BookOfMormon         *Book
	DoctrineAndCovenants *Book
	NewTestament         *Book
	OldTestament         *Book
	PearlOfGreatPrice    *Book
}

type Book struct {
	Books []*SubBook `json:"books"`
}

type SubBook struct {
	Book     string     `json:"book"`
	Chapters []*Chapter `json:"chapters"`
	Sections []*Chapter `json:"sections"`
}

type Chapter struct {
	Chapter   int      `json:"chapter"`
	Section   int      `json:"section"`
	Reference string   `json:"reference"`
	Verses    []*Verse `json:"verses"`
}

type Verse struct {
	Reference string `json:"reference"`
	Text      string `json:"text"`
	Verse     int    `json:"verse"`
}

func (b *Book) Len() int {
	return len(b.Books)
}

func (b *SubBook) Len() int {
	if l := len(b.Chapters); l > 0 {
		return l
	}
	return len(b.Sections)
}

func (c *Chapter) Len() int {
	return len(c.Verses)
}

func (b *Book) GetSubBook(title string) *SubBook {
	for _, book := range b.Books {
		if book.Book == title {
			return book
		}
	}

	return &SubBook{}
}
func (b *Book) GetSubBookN(n int) *SubBook {
	if n >= 0 && n < len(b.Books) {
		return b.Books[n]
	}
	return &SubBook{}
}

func (sb *SubBook) GetChapter(chapter int) *Chapter {
	for _, c := range sb.Chapters {
		if c.Chapter == chapter {
			return c
		}
	}
	for _, c := range sb.Sections {
		if c.Section == chapter {
			return c
		}
	}

	return &Chapter{}
}

func (c *Chapter) GetVerse(verse int) *Verse {
	for _, v := range c.Verses {
		if v.Verse == verse {
			return v
		}
	}

	return &Verse{}
}

func (c *Chapter) Print() {
	fmt.Println(c.Reference)

	for _, v := range c.Verses {
		v.Print()
	}
}

func (v *Verse) Print() {
	fmt.Println(v.Reference)
	fmt.Printf("%d  %s\n", v.Verse, v.Text)
}

func Download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("%v", err)

		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func GetBook(uri string) (*Book, error) {
	// First check if it's already on hdd
	u, _ := url.Parse(uri)
	ps := path.Base(u.Path)

	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	p := path.Join(usr.HomeDir, ".scripture", ps)

	content, err := ioutil.ReadFile(p)
	if err != nil {
		content, err = Download(uri)
		if err != nil {
			return nil, fmt.Errorf("Could not download %s", uri)
		}

		err = ioutil.WriteFile(p, content, 0644)
		if err != nil {
			log.Printf("Was not able to write file to %s: %v", p, err)
		}
	}

	book := &Book{}

	err = json.Unmarshal(content, book)
	if err != nil {
		return nil, err
	}

	if book.Books == nil {
		sub := &SubBook{}

		err = json.Unmarshal(content, sub)
		if err != nil {
			return nil, err
		}

		book.Books = []*SubBook{sub}
	}

	return book, nil
}

func (s Scriptures) GetRandomVerse() *Verse {
	var book *Book

	switch rand.Intn(5) {
	case 0:
		book = s.BookOfMormon
	case 1:
		book = s.DoctrineAndCovenants
	case 2:
		book = s.NewTestament
	case 3:
		book = s.OldTestament
	case 4:
		book = s.PearlOfGreatPrice
	}

	b := rand.Intn(book.Len())
	sub := book.GetSubBook(b)
	c := rand.Intn(sub.Len())
	chapter := sub.GetChapter(c)
	v := rand.Intn(chapter.Len())
	verse := chapter.GetVerse(v)

	return verse
}

func FailIf(err error) {
	if err != nil {
		log.Fatalf("Error encountered: %v", err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var err error

	scriptures := &Scriptures{}

	scriptures.BookOfMormon, err = GetBook(BookOfMormonURL)
	FailIf(err)
	scriptures.DoctrineAndCovenants, err = GetBook(DoctrineAndCovenantsURL)
	FailIf(err)
	scriptures.NewTestament, err = GetBook(NewTestamentURL)
	FailIf(err)
	scriptures.OldTestament, err = GetBook(OldTestamentURL)
	FailIf(err)
	scriptures.PearlOfGreatPrice, err = GetBook(PearlOfGreatPriceURL)
	FailIf(err)

	scriptures.GetRandomVerse().Print()
}
