# About `me3`
My objective with this `me3` repo is to establish a sort of "_mise en place_" to cap the cognitive costs of switching my focus between different personal coding projects.

I would like to have something that resembles some aspects of [Google3](https://dl.acm.org/doi/10.1145/2854146), but for ...just me. That's why it's called `me3`.

## Motivation

This is an experiment. My hypothesis is that the extra effort of maintaining a monorepo (vs multiple distinct but interdependent repos) is justified by the benefits, even in a single-developer use case.

My intuition around this bet is that "Dependency-related tech debt accrued by one dev over a long enough time horizon can be as bad as the dependency-related tech debt accrued by multiple developers over a short time horizon".

## Requirements

- **Support dependency management regardless of source language or target environment**. My personal projects use a variety of different languages, deployment scenarios, and platforms (...sets of external dependencies, etc.).
- **Provide a uniform CLI to build/test/run things**. I would like to use the same set of commands for building, testing and running the code across all of my personal projects.
    - The variety of languages and external dependencies that my projects use naturally results in a set of idiosyncratic build/test/run steps that varies from project to project.
    - Having to re-learn how to build something after not touching it for a long time means that I am less likely to even try, which leads to neglect and bitrot.
- **Repeatable builds** I would like to be able to clone and execute code from my personal monorepo on a new machine or dev environment, and have it "just work" when I want to build/test/execute a fresh checkout in that new environment.
- **Centralize external dependencies**: Installing an external package dependency and getting it to work for one of my projects can take a lot of effort the first time around.  I would like to amortize that effort across all the other projects I have that would also require that external dependency.
- **Work with CI/CD**: Take the path of least resistance here, which in my case is "Work with github actions".

Languages I would like to support in this repo include:
- go
- typescript
- python
- protocol buffers

## Solution

- TL;DR: Bazel

## Inspiration

The initial setup of this repo is *heavily* influenced by things I've learned while working in the `skia-buildbot` repo ([public github mirror of it here](https://github.com/google/skia-buildbot)), so I owe the Skia infra team a healthy debt of gratitude for showing me how it can be done.
