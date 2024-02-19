#include <gtest/gtest.h>
#include "cpp/lib/word_generator.h"

namespace sandbox {
TEST(hello_world_test, hello) {
  auto&& generator = sandbox::create_generator();
  EXPECT_EQ(generator->next(), "Hello, world!");
};
}
