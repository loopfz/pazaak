# pazaak

## What is pazaak ?

Pazaak is a game similar to blackjack that was originally a mini-game for the Star Wars: KOTOR series of videogames.

It is played 1v1, with players taking turns drawing cards ranging from 1 to 10, trying to get the closest to 20 without going over.
Each turn, a player may play a modifier card (+2, -4, ...) from his hand to alter the total.
A game is played as best-of-5 (3 winning), and the hand cards are dealt only once at the start of the match.

[More information on pazaak](http://starwars.wikia.com/wiki/Pazaak/Legends)

## What's in this repository ?

### Image files

There are various apps available, but there are no official physical versions of pazaak cards.
I created image files depicting the card fronts and backs, faithful to the KOTOR look  & feel, which you can use as resources.
These are available in the *Releases* section.

I am not affiliated to BioWare / LucasArts in any way and am only distributing these in an attempt to spread the pazaak love.
If this represents a legal issue I will remove them.

### Command-line game engine

You will find a CLI game engine that can run games and enforce rules.

You can plug bot agents, or possibly human interfaces (untested).

This is useful to simulate many games in a monte-carlo fashion, to analyze the value of modifier cards (+2, -4, +-3, ...).
It may also be used as a basis to develop bot players. A naive but functional implementation is provided in **testagent/**.

To run:

```shell
$ cd pazaakcli
$ go build
$ cd testagent
$ go build
$ cd ..
# stdin line1 = side deck for player 1
# stdin line2 = side deck for player 2
$ ./pazaakcli -player ./testagent/testagent -player ./testagent/testagent <<EOF
auto
auto
EOF
```

### Statistical analysis

The conclusions drawn from analysis conducted using the command-line engine will be provided in the form of a report file.

See **side_deck_card_rank.txt** for a rough early analysis of side deck card value. More to come.
