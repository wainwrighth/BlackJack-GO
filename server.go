package main

import (
	"net/http"
	"math/rand"
	"time"
	"html/template"
	"strconv"
)


// Global variables
var cards [52]int // main deck
var userCards []string
var dealerCards []string
var cc int // index into card array during play
var dealPressed int = 0   // deal button has been pressed
var userTotal int = 0 // User total
var dealerTotal int = 0 // Dealer total
var winTotal int = 0 // Win total
var lossTotal int = 0 // Loss total
var gameComplete bool = false // Keeps track of whether game is over or not
var message string

type Option struct {
    Title string
}

type PageData struct {
    PageTitle string
	Options  []Option
	WinCount string
	LossCount string
	UserHand string
	DealerHand string
	UserCards []string
	DealerCards []string
	Message string
}

// Template to serve for user interface
var tmpl = template.Must(template.ParseFiles("view.html"))
var data = PageData{
	PageTitle: "BlackJack",
	Options: []Option {
		{Title: "new"},
		{Title: "hit"},
		{Title: "stay"},
		{Title: "show"},
		{Title: "shuffle"},
	},
	WinCount: strconv.Itoa(winTotal),
	LossCount: strconv.Itoa(lossTotal),
	UserHand: strconv.Itoa(userTotal),
	DealerHand: strconv.Itoa(dealerTotal),
	UserCards: userCards,
	DealerCards: dealerCards,
	Message: message,
}

// Establish deck and properly shuffle it
func shuffle() {

	// Prevent shuffling if game is currently in play
	if dealPressed != 0 {
		message = "Cannot shuffle deck while hand is in play"
		return
	}

	// Initialize cards array with all 52 cards
	for  i := 0; i < 52; i++ {
		cards[i] = i+1
	}

	// Shuffle cards using rand.Shuffle function
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	// Initialize current position in deck to draw card
	cc = 0
}

// Pull card from the top of the shuffled deck
func getcard() int {
    var cardPull int
	cardPull = cards[cc]
	cc += 1
    return cardPull
}

// Calculate the value of the card based on its index value
func calculate(card int) int {

	if card <= 4 {    // ace
		return 11
	}
	if card <= 20 {   // face card or 10
		return 10
	}
	if card <= 24 {  // 9, and so forth
		return 9
	}    
	if card <= 28 {
		return 8
	}      
	if card <= 32 {
		return 7
	}    
	if card <= 36 {
		return 6
	}    
	if card <= 40 {
		return 5
	}    
	if card <= 44 {
		return 4
	}    
	if card <= 48 {
		return 3
	}
	// Last value left
	return 2
}

func deal() {

	var dlcard1 int
	var dlcard2 int
	var mycard1 int
	var mycard2 int
	dealPressed = 0
	userTotal = 0
	dealerTotal = 0

	// Ensure deck is shuffled before cards are dealt
	if (cc == 0 || cc > 41) {
		shuffle()
	}

	dlcard1 = cards[cc]
	cc++
	dlcard2 = cards[cc]
	cc++
	mycard1 = cards[cc]
	cc++
	mycard2 = cards[cc]
	cc++

	// Update list of card images for display
	userCards = append(userCards, "assets/" + strconv.Itoa(mycard1) + ".png")
	userCards = append(userCards, "assets/" + strconv.Itoa(mycard2) + ".png")
	dealerCards = append(dealerCards, "assets/" + strconv.Itoa(dlcard1) + ".png")
	dealerCards = append(dealerCards, "assets/" + strconv.Itoa(dlcard2) + ".png")
	
	userTotal += calculate(mycard1)
	userTotal += calculate(mycard2)
	dealerTotal += calculate(dlcard1)
	dealerTotal += calculate(dlcard2)
	dealPressed++

    // If the value added to user is an ace, further investigate total
	if (mycard1 >= 1 && mycard1 <= 4) || (mycard2 >= 1 && mycard2 <= 4) {
		// If values added make the user total go over 21 and there is an ace, take off 10 to make the ace worth 1 not 11
		if (userTotal > 21) {
			userTotal = userTotal - 10
		}
	}
	
	// If the value added to dealer is an ace, further investigate total
	if (dlcard1 >= 1 && dlcard1 <= 4) || (dlcard2 >= 1 && dlcard2 <= 4) {
		// If values added make the dealer total go over 21 and there is an ace, take off 10 to make the ace worth 1 not 11
		if (dealerTotal > 21) {
			dealerTotal = dealerTotal - 10
		}
    }
}

func stay() {

	// Game is already finished, do nothing
	if gameComplete {
		message = "Game is complete, can't stay"
	// If deal has been initiated, settle dealer score and determine outcome
	} else if dealPressed != 0 {
		dealPressed = 0

		// User busted, record loss
		if userTotal > 21 {
			lossTotal++
			gameComplete = true
			message = "Dealer wins :("
			return
		}

		// If dealer has under 17, they should draw
		if dealerTotal < 17 {
			for {
				if dealerTotal > 17 {
					break
				}
				dealerCard := getcard()
				cardVal := calculate(dealerCard)
				dealerTotal += cardVal
				dealerCards = append(dealerCards, "assets/" + strconv.Itoa(dealerCard) + ".png")

				// Check if the card drawn is an ace and count it based on total
				if dealerCard >= 1 && dealerCard <= 4 {
					if (dealerTotal > 21) {
						dealerTotal = dealerTotal - 10
					}
				}
			}
		}

		// If dealer busted, record win
		if dealerTotal > 21 {
			winTotal++
			gameComplete = true
			message = "User wins!"
		}

		// if dealer ends in range past 17 and no bust determine outcome
		if dealerTotal >= 17 && dealerTotal <= 21 {
			if dealerTotal > userTotal {
				lossTotal++
				gameComplete = true
				message = "Dealer wins :("
			} else if userTotal > dealerTotal  || userTotal == dealerTotal {
				winTotal++
				gameComplete = true
				message = "User wins!"
			}
		}
	}
}

func updatePage() {
	data.DealerHand = strconv.Itoa(dealerTotal)
	data.UserHand = strconv.Itoa(userTotal)
	data.WinCount = strconv.Itoa(winTotal)
	data.LossCount = strconv.Itoa(lossTotal)
	data.DealerCards = dealerCards
	data.UserCards = userCards
	data.Message = message
}

func newHandler(w http.ResponseWriter, r *http.Request) {

	// Reset Values
	dealerCards = nil
	userCards = nil
	gameComplete = false
	message = ""

	deal()
	updatePage()
	tmpl.Execute(w, data)
}

func hitHandler(w http.ResponseWriter, r *http.Request) {

	if userTotal > 21 {
		message = "User busted, hit 'stay' to end hand"
	}

	// Deal hasnt been initiated, dont take card
	if dealPressed != 0  && userTotal <= 21 {
		card := getcard()
		val := calculate(card)
		userTotal += val
		userCards = append(userCards, "assets/" + strconv.Itoa(card) + ".png")

		// If card is an ace, count its value according to hand
		if (card >= 1 && card <= 4) || (card >= 1 && card <= 4) {
			if (userTotal > 21) {
				userTotal = userTotal - 10
			}
		}

		updatePage()
		tmpl.Execute(w, data)
	} else { // Do nothing
		updatePage()
		tmpl.Execute(w, data)
	}
}

func stayHandler(w http.ResponseWriter, r *http.Request) {
	stay()
	updatePage()
	tmpl.Execute(w, data)
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	updatePage()
	tmpl.Execute(w, data)
}

func shuffleHandler(w http.ResponseWriter, r *http.Request) {
	shuffle()
	updatePage()
	tmpl.Execute(w, data)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	updatePage()
	tmpl.Execute(w, data)
}

func main(){
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/new", newHandler)
	mux.HandleFunc("/hit", hitHandler)
	mux.HandleFunc("/stay", stayHandler)
	mux.HandleFunc("/show", showHandler)
	mux.HandleFunc("/shuffle", shuffleHandler)
	mux.HandleFunc("/view", viewHandler)
	http.ListenAndServe(":80", mux)
}
