import os
import random

# TODO: Move this into an init function of some kind.
ws = os.path.dirname(__file__)
wordsfile = open(os.path.join(ws, "wordlist.txt"), 'r')
words = wordsfile.readlines()
words = list(map(lambda x: x.strip(), words))

# TODO: Add a seed parameter to this function for the rng.
def new_board(words, first_team="blue"):
  bag = list(words)
  random.shuffle(bag)

  idx = 8 if first_team == "blue" else 9
  red = bag[:idx]
  bag = bag[idx:]

  idx = 8 if first_team == "red" else 9
  blue = bag[:idx]
  bag = bag[idx:]

  bystanders = bag[:7]
  bag = bag[7:]

  assassin = bag[0]

  board = {
      'red': red,
      'blue': blue,
      'bystanders': bystanders,
      'assassin': assassin
  }
  return board
