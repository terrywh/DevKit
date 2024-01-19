#!/bin/bash

echo "upgrade kubectl:"
deno run --allow-net --allow-read --allow-run --allow-write script/kubectl.js

echo "upgrade trzsz:"
deno run --allow-net --allow-read --allow-run --allow-write script/trzsz.js