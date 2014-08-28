package hello

import (
  "net/http"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/cors"
  "github.com/martini-contrib/render"
  "math/rand"
  "strconv"
  "sort"
)

 type Resources struct {
	Name  string `json:"name"`
	URI  string `json:"uri"`
	Methods  string `json:"methods"`
 }

 type Card struct {
	Suit  string `json:"suit"`
	Number int `json:"number"`
 }

//deck
var cards = []Card{{Suit: "Spades", Number: 2},{Suit: "Spades", Number: 3},{Suit: "Spades", Number: 4},{Suit: "Spades", Number: 5},{Suit: "Spades", Number: 6},{Suit: "Spades", Number: 7},{Suit: "Spades", Number: 8},{Suit: "Spades", Number: 9},{Suit: "Spades", Number: 10},{Suit: "Spades", Number: 11},{Suit: "Spades", Number: 12},{Suit: "Spades", Number: 13},{Suit: "Spades", Number: 14},{Suit: "Hearts", Number: 2},{Suit: "Hearts", Number: 3},{Suit: "Hearts", Number: 4},{Suit: "Hearts", Number: 5},{Suit: "Hearts", Number: 6},{Suit: "Hearts", Number: 7},{Suit: "Hearts", Number: 8},{Suit: "Hearts", Number: 9},{Suit: "Hearts", Number: 10},{Suit: "Hearts", Number: 11},{Suit: "Hearts", Number: 12},{Suit: "Hearts", Number: 13},{Suit: "Hearts", Number: 14},{Suit: "Diamonds", Number: 2},{Suit: "Diamonds", Number: 3},{Suit: "Diamonds", Number: 4},{Suit: "Diamonds", Number: 5},{Suit: "Diamonds", Number: 6},{Suit: "Diamonds", Number: 7},{Suit: "Diamonds", Number: 8},{Suit: "Diamonds", Number: 9},{Suit: "Diamonds", Number: 10},{Suit: "Diamonds", Number: 11},{Suit: "Diamonds", Number: 12},{Suit: "Diamonds", Number: 13},{Suit: "Diamonds", Number: 14},{Suit: "Clubs", Number: 2},{Suit: "Clubs", Number: 3},{Suit: "Clubs", Number: 4},{Suit: "Clubs", Number: 5},{Suit: "Clubs", Number: 6},{Suit: "Clubs", Number: 7},{Suit: "Clubs", Number: 8},{Suit: "Clubs", Number: 9},{Suit: "Clubs", Number: 10},{Suit: "Clubs", Number: 11},{Suit: "Clubs", Number: 12},{Suit: "Clubs", Number: 13},{Suit: "Clubs", Number: 14}}

//shuffle
func Shuffle(slc []Card) {
    for i := 1; i < len(slc); i++ {
        r := rand.Intn(i + 1)
        if i != r {
            slc[r], slc[i] = slc[i], slc[r]
        }
    }
}





// mutli sorting taken from http://golang.org/pkg/sort/

type lessFunc func(p1, p2 *Card) bool

// multiSorter implements the Sort interface, sorting the cards within.
type multiSorter struct {
	cards []Card
	less    []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiSorter) Sort(cards []Card) {
	ms.cards = cards
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}


func (ms *multiSorter) Len() int {
	return len(ms.cards)
}

func (ms *multiSorter) Swap(i, j int) {
	ms.cards[i], ms.cards[j] = ms.cards[j], ms.cards[i]
}

func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.cards[i], &ms.cards[j]
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
	}
	return ms.less[k](p, q)
}






func init() {
  m := martini.Classic()
  
  allowCORSHandler := cors.Allow(&cors.Options{
	AllowOrigins:     []string{"http://*.techslides.com"},
	AllowMethods:     []string{"GET", "POST"},
	AllowHeaders:     []string{"Origin"},
  })

m.Use(render.Renderer(render.Options{
  IndentJSON: true, // Output human readable JSON
}))


  m.Get("/", allowCORSHandler, func(r render.Render) {
  	Shuffle(cards)
	json := []Resources{{Name: "cards", URI: "/cards/{players}", Methods: "GET"}}
	r.JSON(200, json)
  })

  m.Get("/cards", allowCORSHandler, func(r render.Render) {
  	Shuffle(cards)
	r.JSON(200, cards)
  })  

  m.Get("/cards/:players", allowCORSHandler, func(args martini.Params, r render.Render) {
  	
  	Shuffle(cards)
  	
  	//use Atoi as ParseInt always does int64
  	p, _ := strconv.Atoi(args["players"])

	if(p==0 || p>52){
		//not valid number
		r.JSON(200, cards)
    } else {

		var total = len(cards)
		var start = 0
    	var split = total/p

		//for sorting by suit and number
	    suit := func(c1, c2 *Card) bool {
			return c1.Suit < c2.Suit
		}
		numbers := func(c1, c2 *Card) bool {
			return c1.Number < c2.Number
		}

    	//make a multidiemnsional array in Go (http://rosettacode.org/wiki/Create_a_two-dimensional_array_at_runtime#Go) with as many items in array as players so we can loop and fill below
    	var a = make([][]Card,p)

	    for i := range a {
	    	var psplit = cards[start:start+split];

	    	// Sort by Number and Suit
			OrderedBy(suit, numbers).Sort(psplit)

	    	a[i]=psplit
	    	start=start+split

	    }

	    //kitty
	    var kittycards = cards[start:total]

    	r.JSON(200, map[string]interface{}{"players": a, "kitty":kittycards})
    }


  })  
  http.Handle("/", m)
}