## findsubstring

### Usage

You can get usage information like this:

    findsubstrings -help

Which displays:

    usage: findsubstrings [options]
      -help
            show this help message
      -k int
            minimum substring length, used for finding match seeds (default 8)
      -patterns string
            filename containing substrings to look for (required)


### Example

Suppose we have the following input

    hello there
    bellow
    what is that hello
    not a match here
    what the hell
    two hellos in one line hello again

And we have the following patterns file:

    hello
    hell
    ellow

We could run the program as follows:

    cat unit_input | findsubstrings -k 4 -patterns unit_patterns > unit_out

This would display the following to standard error:

    stored 5 seeds
    found 5 matches
    hello 3
    hell 2
    ellow 0

And the `unit_out` file would contain:

    0	0	hell	hello there
    0	0	hello	hello there
    5	4	hell	two hellos in one line hello again
    5	4	hello	two hellos in one line hello again
    5	23	hello	two hellos in one line hello again

