# BlackJack

## Files

server.go - main go file with web server and blackjack game endpoints in it

view.html - html file with the template the web server uses to display the game to the user

/assets/ - directory with card images, css style sheet and favicon in it

## Steps

Once the BlackJack directory is unzipped the server can be ran using (The server defaults to port 80):

>$go run server.go

Once the server is running open a browser and navigate to:

>$localhost:80/view

This is the main page for the BlackJack game displaying the gaming options to the user and stats

## Notes

- All endpoints specified in the google doc are accessible through their respective endpoints, some just won't display information to the user because the format the page is at doesn't cause there to be a need (example: localhost:80/show)
- After picking to stay, the dealer will hit only if they are under to 17 becuase only then is it in their favor to do so. After the dealer is finished the final result is decided. In the event the user and dealer tie, the user is awarded a win
- The color of the top buttons is in fact rebeccapurple
