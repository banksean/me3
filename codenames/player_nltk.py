import nltk
from nltk.stem import PorterStemmer
from nltk.corpus import wordnet as wordnet_corpus
from powerset import powerset

ps = PorterStemmer()


def related_words(wn, word):
    hypo = lambda s: s.hyponyms()
    hyper = lambda s: s.hypernyms()
    ret = set()
    synonyms = []
    antonyms = []

    # More specific terms
    hyponyms = []

    # More broad terms
    hypernyms = []

    for s in wn.synsets(ps.stem(word)):
        hyponyms.extend([s.lemmas()[0].name() for s in s.closure(hypo, depth=1)])
        hypernyms.extend([s.lemmas()[0].name() for s in s.closure(hyper, depth=1)])

        for l in s.lemmas():
            synonyms.append(l.name())
            for a in l.antonyms():
                antonyms.append(a.name())

    return {
        "synonyms": list(set(synonyms)),
        "antonyms": list(set(antonyms)),
        "hyponyms": list(set(hyponyms)),
        "hypernyms": list(set(hypernyms)),
    }


def wordnet_hints(wn, board, player="blue"):
    opponent = "red" if player == "blue" else "red"

    # neg_terms should contain all of the opponent's codewords
    # plus all of those codewords' related words (so hints don't point to your
    # opponent's codewords).
    neg_terms = set()
    for word in board[opponent]:
        word = ps.stem(word)
        rel_words = related_words(wn, word)
        neg_terms.union(set(rel_words["synonyms"]))
        neg_terms.union(set(rel_words["antonyms"]))
        neg_terms.union(set(rel_words["hyponyms"]))
        neg_terms.union(set(rel_words["hypernyms"]))

    # iterate over 1-3 word permutations of board[player]
    p = powerset(board[player])
    hints = []

    for s in p:
        if len(s) == 0:
            continue
        if len(s) > 3:
            break

        pos_terms = []
        for w in s:
            w = ps.stem(w)
            rel_words = related_words(wn, w)
            all_related = set(rel_words["synonyms"])
            all_related.union(set(rel_words["antonyms"]))
            all_related.union(set(rel_words["hyponyms"]))
            all_related.union(set(rel_words["hypernyms"]))
            pos_terms.append(all_related)

        pos_terms = pos_terms[0].intersection(*pos_terms)

        candidates = list(pos_terms.difference(neg_terms))
        if len(candidates) == 0:
            continue

        # print(s)
        # print(" candidate hints: ", candidates)
        # print(" conflicting hints: ", neg_terms.intersection(pos_terms))

        for term in candidates:
            if (
                term.lower() in pos_terms | neg_terms | set(s)
                or term.lower() is board["assassin"]
                or "_" in term
            ):
                continue
            hints.append({"targets": s, "hint": term})

    hints = sorted(hints, key=lambda x: len(x["targets"]), reverse=True)
    return hints


def wordnet_guess(wn, clue, n, terms):
    ret = set()
    for i in range(n):
        guesses = set()
        # g = model.most_similar_to_given(clue, list(terms))
        # ret.add(g)
        # terms.discard(g)

        rel_words = related_words(wn, clue)
        # print('related words: ', json.dumps(rel_words, indent=2))
        guesses = guesses.union(set(rel_words["synonyms"]))
        guesses = guesses.union(set(rel_words["antonyms"]))
        guesses = guesses.union(set(rel_words["hyponyms"]))
        guesses = guesses.union(set(rel_words["hypernyms"]))

        # print('guesses for ', clue, ': ', guesses)
        guesses = guesses.intersection(terms)
        # print('guesses filtered for ', clue, ': ', guesses)
        if len(guesses) > 0:
            g = list(guesses)[0]
            ret.add(g)
            terms.discard(g)

    return ret


def load_wordnet():
    # assert(nltk.download('wordnet'))  # Make sure we have the wordnet data.
    return wordnet_corpus
