# Exercise 03: Choose Your Own Adventure

Gets a JSON file, throught flag configuration and shows its information as an interactive story.

### JSON Structure
``` json
{
    # Arc Name
    "intro": {
        "title": string,
        "story": [ string ],
        "options": [
        {
            "text": string,
            "arc": string
        },
        {
            "text": string,
            "arc": string
        }
        ]
    }
}

```

### Functional Options

Some information about this pattern used in this exercise.

* [Functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
* [go-patterns: Functional Options](https://github.com/tmrts/go-patterns/blob/master/idiom/functional-options.md)