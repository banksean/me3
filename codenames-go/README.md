# An LLM-Powered "Codenames" Player

## TL;DR:
Terminal-based Codenames player that can use LLMs to play the role of any or all of the participants in the game.

### Example: LLM (Field Agent and Spy Master) vs LLM (Field Agent and Spy Master)

To watch Llama3 play against Gemma:
```
bazel run //codenames-go -- \
  -red-spy-master=llama3 \
  -red-field-agent=llama3 \
  -blue-spy-master=gemma \
  -blue-field-agent=gemma
```

[![asciicast](https://asciinema.org/a/Tr7xa9VahwIbl8SpzMDrfmUMv.svg)](https://asciinema.org/a/Tr7xa9VahwIbl8SpzMDrfmUMv)

### Example: Human Field Agent + LLM Spy Master vs LLM (Field Agent and Spy Master)

To play on the red team as a human field agent, with Gemma as your spymaster, against a blue team with Llamma3 playing both spymaster and field agent:

```
bazel run //codenames-go -- \
  -red-spy-master=gemma \
  -red-field-agent=interactive \
  -blue-spy-master=llama3 \
  -blue-field-agent=llama3
```
[![asciicast](https://asciinema.org/a/JYv6z2opEasDkEJE8WVj5qIxv.svg)](https://asciinema.org/a/JYv6z2opEasDkEJE8WVj5qIxv)


## Background
### Codenames, the game
[Codenames](https://en.wikipedia.org/wiki/Codenames_(board_game)) is a popular tabletop game played by two teams (Red, and Blue). Each team has one player (Spy Master) who provides clues to all the other players (Field Agents) on their team.

[BoardGameGeek](https://boardgamegeek.com/boardgame/178900/codenames) describes it well: 

> Two rival spymasters know the secret identities of 25 agents. Their teammates know the agents only by their 
codenames — single-word labels like "disease", "Germany", and "carrot". Yes, carrot. It's a legitimate codename. Each spymaster wants their team to identify their agents first...without uncovering the assassin by mistake.
>
> In Codenames, two teams compete to see who can make contact with all of their agents first. Lay out 25 cards, each bearing a single word. The spymasters look at a card showing the identity of each card, then take turns clueing their teammates. A clue consists of a single word and a number, with the number suggesting how many cards in play have some association to the given clue word. The teammates then identify one agent they think is on their team; if they're correct, they can keep guessing up to the stated number of times; if the agent belongs to the opposing team or is an innocent bystander, the team's turn ends; and if they fingered the assassin, they lose the game.
>
> Spymasters continue giving clues until one team has identified all of their agents or the assassin has removed one team from play.

### Writing automated teammate/opponent players for Codenames

A few years ago I attempted to build an automated Codenames players using more primitive approaches that were available at the time, like [Word2Vec](https://en.wikipedia.org/wiki/Word2vec) and [NLTK's wordnet](https://www.nltk.org/howto/wordnet.html) package. These did not work terribly well in my opinion, so I shelved the idea.

Fast-forward to the recent rise in popularity and accessibility of LLMs, I had to dust off the idea and try it again.


## How this code automates player roles for the game

### Specify who is playing which role on each team
Since players can be on one of two teams, each player can play one of two roles on each team, the four team/role combinations are specified on the command line.

```
-red-spy-master=<Player Type> 
-red-field-agent=<Player Type>
-blue-spy-master=<Player Type> 
-blue-field-agent=<Player Type>
```

Each team/role flag value can be either "`interactive`" (the default, which blocks on terminal input from a human user) or the name of a model supported by Ollama ("`llama3`", "`gemma`" etc.).

By convention, RED team always goes first.

### Gameplay

If any of the roles are `interactive`, the application will block on terminal input from the user, after printing a prompt describing which team and role should provide the input.

Otherwise, the turns are executed as fast as the LLMs respond.  Playing two LLMs against each other locally can be slower that playing the same LLM against itself due to how OLlama 

### Game State and LLM interaction

All the game state (what cards are in play, which cards have been guessed/revealed so far, which cards belong to which team, whose turn it currently is, etc.) is managed by handwritten Go code. It treats the LLM as a fairly dumb, stateless service for generating clues and making guesses based on clues.

The prompts do not assume any ongoing context to reflect the current state of the game, so all relevant information is provided to the LLM in each query.  I.e. the LLM doesn't know it's playing a game called "codenames" or that there are teams and it is playing the role of a player on one of those teams, etc.

The prompts are essentially of the form, "What's a word that hints at one of these words [...] but NOT these words [...]?", or "Based on this hint: '...', what is the best matching word from this list: [...]?"

There are of course much more advanced approaches one could use that treat the LLM as a smarter participant, and I came to the above approach after trying some of them. 

For example, I tried using ChatGPT to implement a Codenames player by giving it only what a human would get: the copied and pasted text of the game's instruction manual.  What I found was that it sort of understood some of the features of the game, but it was very unreliable when it came to know whose turn it was, which cards have been revealed already, the secret identity of each card in play etc.  
