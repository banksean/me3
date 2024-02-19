#include <iostream>
#include <string>
#include <vector>

#include "absl/strings/str_join.h"

#include "word_generator.h"

int main() {
  std::vector<std::string> v = {"foo", "bar", "baz"};
  std::string s = absl::StrJoin(v, "-");

  std::cout << "Joined string: " << s << "\n";

  auto&& generator = sandbox::create_generator();
  std::cout << generator->next() << std::endl;

  return 0;
}