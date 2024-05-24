import random
from board import new_board, words
from player_word2vec import load_word2vec
from player_nltk import load_wordnet, wordnet_hints, wordnet_guess


def evaluate_wordnet(wn, seed):
    random.seed(seed)
    board = new_board(words)
    hint_list = wordnet_hints(wn, board, player="blue")

    # red_hints = wordnet_hints(wn, board, player='red')

    all_words = set(board["red"]) | set(board["blue"]) | set([board["assassin"]])
    # print(hint_list)
    print(len(hint_list), " hints")

    score = 0
    total_possible = 0
    for n in range(len(hint_list)):
        g = wordnet_guess(
            wn,
            hint_list[n]["hint"],
            len(hint_list[n]["targets"]),
            all_words,  # | set(board['assassin']) # Whole board
        )
        correct = set(board["blue"]) & set(g)
        incorrect = set(board["red"]) & set(g)
        total_possible += len(hint_list[n]["targets"])

        # TODO: actual scoring. But for now:
        # Ignore target words from the hints, just look at hint, n and what's on the board.
        # If all n of the guesses are on your squares (regardless of target words) then
        # that's n points.
        score += len(correct) - len(incorrect)

        # print('  hint: ', hint_list[n])
        # print('  guess: ', g)
        # print('correct: ', correct)

    if total_possible > 0:
        print("total score: ", score, " ", (score / total_possible))
    else:
        print("could not generate any hints!")
    # print(json.dumps(board, indent=2))

    # print(json.dumps(blue_hints, indent=2))
    # print(json.dumps(red_hints, indent=2))


if __name__ == "__main__":
    board = new_board(words)
    print(board)
    # w2vmodel = load_word2vec()
    # print(w2vmodel)
    wnmodel = load_wordnet()
    print(wnmodel)
    for s in range(1):
        print("seed: ", s)
        evaluate_wordnet(wnmodel, s)
