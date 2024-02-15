import time

import gensim.downloader as api
from powerset import powerset

RESTRICT_VOCAB=5000
SCORE_THRESHOLD=0.0 #1936

def load_word2vec():
    print("loading word2vec-google-news-300\n")
    t1 = time.perf_counter(), time.process_time()
    model = api.load('word2vec-google-news-300')
    t2 = time.perf_counter(), time.process_time()
    
    print("word2vec-google-news-300 loaded\n")
    print(f" Real time: {t2[0] - t1[0]:.2f} seconds")
    print(f" CPU time: {t2[1] - t1[1]:.2f} seconds")

    print("most similar to 'cat': %s\n" % model.most_similar("cat")) # show words that similar to word 'cat'
    t3 = time.perf_counter(), time.process_time()
    print(f" Real time: {t3[0] - t2[0]:.2f} seconds")
    print(f" CPU time: {t3[1] - t2[1]:.2f} seconds")

    return model

def word2vec_guess(model, clue, n, terms):
  ret = set()
  terms = set([term for term in terms if term in model.vocab])
  for i in range(n):
    g = model.most_similar_to_given(clue, list(terms))
    ret.add(g)
    terms.discard(g)
  return ret

# Make sure to add the asassin to negative (the other team's terms) before calling this.
def word2vec_hints(model, board, player='blue'):
  opponent = 'red' if player == 'blue' else 'red'
  positive = set([w for w in board[player] if w in model.vocab])
  negative = set([w for w in board[opponent] if w in model.vocab])
  if board['assassin'] in model.vocab:
    negative.add(board['assassin'])

  # TODO: Expand negative with N similar terms to each word in negative?
  # This may already happen due to vector math in word2vec.

  # Enumerate combinations of 1, 2, 3 terms
  p = powerset(positive)
  clues = []
  for s in p:
    if len(s) == 0:
      continue
    if len(s) > 3:
      break
    candidates = model.most_similar(
      list(s),
      list(negative),
      topn=10,
      restrict_vocab=RESTRICT_VOCAB
    )
    hi_score = 0.0
    hi_clue = None

    for term, score in candidates:
      if term.lower() in positive | negative or term.startswith('afp'):
        continue
      if score > hi_score and score > SCORE_THRESHOLD:
        hi_score = score
        hi_clue = term

    if hi_clue != None:
      clues.append({'targets': s, 'hint': hi_clue, 'score': hi_score})

  # This ranking penalizes/benefits the higher-term-count clues.
  # TODO: Tune this.
  clues.sort(key=lambda x: x['score']*len(x['targets']), reverse=True)
  return clues
