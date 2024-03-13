# call-in-parallel

```
call-in-parallel - Run a command multiple times in parallel

  Flags:
       --version   Displays the program version string.
    -h --help      Displays help with available flag, subcommand, and positional value parameters.
    -n           number of instances to start (default: 1)
    -d --delay     delay between starting each instance (default: 10ms)
```

## Installation

```
go install github.com/tilebox/call-in-parallel@latest
```

## Usage Examples

Print "hello" 3 times in parallel, with a default delay of 10ms between each command
```bash
> call-in-parallel -n 3 -- echo Hello World!
```

Fetch the [tilebox.com](https://tilebox.com) website 5 times in parallel, with a delay of 1s in-between
```bash
> call-in-parallel -n 5 -d 1s -- curl https://tilebox.com
```
