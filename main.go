package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
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

// If there are any errors book will be nil.
func DownloadBook(url string) (book *Book) {
	book = &Book{}

	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)

		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)

		return
	}

	err = json.Unmarshal(b, book)
	if err != nil {
		log.Print(err)

		return
	}

	return
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

	b := rand.Intn(len(book.Books))
	sub := book.Books[b]
	c := rand.Intn(len(sub.Chapters))
	chapter := sub.Chapters[c]
	v := rand.Intn(len(chapter.Verses))
	verse := chapter.Verses[v]

	return verse
}

func main() {
	rand.Seed(time.Now().UnixNano())
	scriptures := &Scriptures{
		BookOfMormon:         DownloadBook(BookOfMormonURL),
		DoctrineAndCovenants: DownloadBook(DoctrineAndCovenantsURL),
		NewTestament:         DownloadBook(NewTestamentURL),
		OldTestament:         DownloadBook(OldTestamentURL),
		PearlOfGreatPrice:    DownloadBook(PearlOfGreatPriceURL),
	}

	scriptures.GetRandomVerse().Print()
}
