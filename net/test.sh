#!/bin/bash

# Define the output file
OUTPUT_FILE="test/test_output.log"
FILTERED_OUTPUT_FILE="test/test_output_filtered.log"

# Run the command and capture its output

echo "Running tests..."
#if something goes wrong use full path on linux /usr/local/share/dotnet/x64/dotnet
dotnet test './test/App.Tests.sln' > "$OUTPUT_FILE" 2>&1

ERROR_CODE=$?

grep -v '/App.Tests/' "$OUTPUT_FILE" > "$FILTERED_OUTPUT_FILE" 2>&1

mv "$FILTERED_OUTPUT_FILE" "$OUTPUT_FILE"

if [ $ERROR_CODE -eq 0 ]; then
  echo "Tests passed"
else
  echo "Tests failed"
fi


cat "$OUTPUT_FILE"
exit $ERROR_CODE