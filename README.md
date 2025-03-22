# Content

- `cms` - terminal interface to interact with content sqlite database.
- `site` - run the static site generator.

## Usage

Build everything with the default justfile recipie:

`make`

## Usage 

`./cms --help`

`./site`

How to: goose migrations

Here are the steps:

[install goose](https://pressly.github.io/goose/installation/)
1. Set environment variables for goose.

`source goose.env`

2. Create a goose sql migration
