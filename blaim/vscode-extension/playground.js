// This is a playground for testing generative code suggestion logging.

function fibonacci(n) {
  if (n < 2) {
    return n;
  } else {
    return fibonacci(n - 1) + fibonacci(n - 2);
  }
}

// This is hand-written:
function testFib() {
  let v = fibonacci(1);
  console.log(v);
}
