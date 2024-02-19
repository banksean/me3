#pragma once

#include <memory>
#include <string>

namespace sandbox {

struct word_generator {
    virtual ~word_generator() = default;

    virtual std::string next() const = 0;

    virtual std::string next(unsigned times) const = 0;
};

std::unique_ptr<word_generator> create_generator();
};