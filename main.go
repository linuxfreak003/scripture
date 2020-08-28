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
}

type Chapter struct {
	Chapter   int      `json:"chapter"`
	Reference string   `json:"reference"`
	Verses    []*Verse `json:"verses"`
}

type Verse struct {
	Reference string `json:"reference"`
	Text      string `json:"text"`
	Verse     int    `json:"verse"`
}

func (b *Book) Get(title string) *SubBook {
	for _, book := range b.Books {
		if book.Book == title {
			return book
		}
	}

	return &SubBook{}
}

func (sb *SubBook) Get(chapter int) *Chapter {
	for _, c := range sb.Chapters {
		if c.Chapter == chapter {
			return c
		}
	}

	return &Chapter{}
}

func (c *Chapter) Get(verse int) *Verse {
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

// If there are any errors book will be nil.
func DownloadBook(url string) (book *Book) {
	book = &Book{}
	b, err := Download(url)
	if err != nil {
		log.Printf("%v", err)

		return
	}
	err = json.Unmarshal(b, book)
	if err != nil {
		log.Printf("%v", err)

		return
	}

	return
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
	return book, nil
}

func (s Scriptures) GetRandomVerse() *Verse {
	var book *Book

	switch rand.Intn(5) {
	case 0:
		book = s.BookOfMormon
	case 1:
		// book = s.DoctrineAndCovenants
		// Broken until the structs can be fixed
		// D&C Scructure is different
		book = s.BookOfMormon
	case 2:
		book = s.NewTestament
	case 3:
		book = s.OldTestament
	case 4:
		book = s.PearlOfGreatPrice
	}

	b := rand.Intn(len(book.Books))
	sub := book.Books[b]
	c := rand.Intn(len(sub.Chapters))
	chapter := sub.Chapters[c]
	v := rand.Intn(len(chapter.Verses))
	verse := chapter.Verses[v]

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
