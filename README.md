# Azul Plugin Entropy

This plugin is responsible for calculating Shannon's entropy across
blocks of a given binary file. This is published as an `info` field
for viewing purposes and an overall entropy added as a feature to the binary.

Shannon's entropy formula:
= - SUM i(p(i)×log2(p(i)))

Where p(i) is the probability a randomly selected value (within the sample) would be that value.

Entropy calculates entropy for an entire file and chunks the file into a minimum of 256byte blocks, and
a maximum of 800 file blocks, and calculates the entropy for each of those blocks.

## Potential for ignoring blocks

Entropy will ignore the last bytes in a file for the chunked file entropies, there may be multiple blocks worth of data
ignored or just a partial block worth.

This is caused by requiring all blocks have the same block_size and having a hard limit on the number of blocks.

A detailed example is here:
For a block size of 800 and a file that has 307,618bytes
The size of each block will be calculated as size=384, and the blocks=800, formula used:
size = contentSize / blocks = 307,618 / 800 = 384.5225
The size is then integer rounded down.

So the total amount of data that fits into the blocks is:
800blocks \* 384bytes = 307,200bytes which with the original file size will have another block 307,200 + 384 < 307,618.
Entropy will ignore this block.

The reason entropy doesn't use a size of 385bytes in this case is that it would result in insufficient data for the
final block as 800blocks \* 385 = 308,000bytes which would mean the last block would have insufficient data because
307,618 - (308,000 - 385) = 3bytes for the final bucket and 256 is the minimum to do an entropy calculation.

## Events

Events Consumed:

- entity_type: `binary`, event: !`binary_enriched`

Events Produced:

- entity_type: `binary`, event: `binary_enriched`

## Usage

    PLUGIN_DATA_URL=http://localhost:8111 PLUGIN_EVENTS_URL=http://localhost:8111 azul-entropy

## Local Build

`go build -v -tags netgo -ldflags '-w -extldflags "-static"' -o bin/azul-entropy *.go`

## Docker Builds

An example dockerfile is provided for building images.
To use the container for a build run the following (or similar if your ssh private and public key for accessing Azure is in a non-standard file):

Example Build (requires you install `buildah` with `sudo apt install buildah`):

```bash
buildah build --volume ~/.ssh/known_hosts:/root/.ssh/known_hosts --ssh id=~/.ssh/id_rsa  .
```
