#!/bin/sh

set -e

###############################################################################
#
# This is the script copied into the docker image for the scan function.  It
# is essentially the entrypoint to the function
#
# Required Environment
# -----------------------
# GITHUB_API_URL: Github API URL
# GITHUB_ORG: Github Organization
# GITHUB_REPO: Github Repository
# GITHUB_REF: Github Ref, i.e - main, master, v1.3, 6faed2
# OUTPUT_FORMAT: Output Format of the scan
#
# Optional Environment
# -----------------------
# GITHUB_AUTH_TOKEN: Github authorization token if required for accessing repo
#
# Optional Standard Input (stdin)
# If the stdin is populated, use it as the scan configuration
#
# Exit Code Handling
# -----------------------
# By default exit codes are delegated to shell and commands.
#
# Custom Exit Codes
# -----------------------
# 61: Github Authorization Failure
# 64: Github Org / Repo / Ref Not Found
# 65: Error Downloading Source Code
#
###############################################################################

UUID=$(date +%s%N | cut -b1-13)
ZIP_FILE="${GITHUB_ORG}-${GITHUB_REPO}-${GITHUB_REF}-$UUID.zip"
SOURCE_CODE_EXTRACT_DIRECTORY="${GITHUB_ORG}-${GITHUB_REPO}-${GITHUB_REF}-$UUID"
PROVIDED_CONFIG_FILE=".sourcehawk-provided.yml"

error_and_exit() {
  echo "$2" > /dev/stderr
  exit "$1"
}

cleanup() {
  rm -rf "$ZIP_FILE" "$SOURCE_CODE_EXTRACT_DIRECTORY" "$PROVIDED_CONFIG_FILE-"
}

trap cleanup INT EXIT

# Retrieve the provided config from stdin
if [ -p /dev/stdin ]; then
  cat > "$PROVIDED_CONFIG_FILE"
fi

# Download a zip file of the repository contents
GITHUB_ZIPBALL_URL="$GITHUB_API_URL/repos/$GITHUB_ORG/$GITHUB_REPO/zipball/$GITHUB_REF"
CURL_RESPONSE_CODE=0
if [ -z "${GITHUB_AUTH_TOKEN}" ]; then
  CURL_RESPONSE_CODE=$(curl -L -s -w "%{http_code}" "$GITHUB_ZIPBALL_URL" -o "$ZIP_FILE" || error_and_exit 65 "Error downloading/writing zip file")
else
  CURL_RESPONSE_CODE=$(curl -L -s -w "%{http_code}" -H "Authorization: token $GITHUB_AUTH_TOKEN" "$GITHUB_ZIPBALL_URL" -o "$ZIP_FILE" || error_and_exit 65 "Error downloading/writing zip file")
fi
if [ "$CURL_RESPONSE_CODE" = 401 ]; then
  if [ -z "${GITHUB_AUTH_TOKEN}" ]; then
    error_and_exit 61 "Are you trying to scan a non-public repository? An authorization token is required."
  else
    error_and_exit 61 "The authorization token provided is not valid"
  fi
elif [ "$CURL_RESPONSE_CODE" = 404 ]; then
  error_and_exit 64 "Could not find source code to scan.  If the repository requires access, make sure to provide an Authorization token"
fi

# Unzip the source code tarball
mkdir -p "$SOURCE_CODE_EXTRACT_DIRECTORY" && unzip -qq -d "$SOURCE_CODE_EXTRACT_DIRECTORY" "$ZIP_FILE"

# Grab the repository root from unzipped tarball
# shellcheck disable=SC2035
SOURCE_CODE_ROOT_DIRECTORY=$(cd "$SOURCE_CODE_EXTRACT_DIRECTORY" && cd */. && pwd)

# Use the provided config file if present, otherwise default
CONFIG_FILE="${SOURCE_CODE_ROOT_DIRECTORY}/sourcehawk.yml"
if [ -f "$PROVIDED_CONFIG_FILE" ]; then
  CONFIG_FILE="$(pwd)/$PROVIDED_CONFIG_FILE"
fi

# Execute the scan
sourcehawk scan --verbosity MEDIUM --output-format "$OUTPUT_FORMAT" --config-file "$CONFIG_FILE" "$SOURCE_CODE_ROOT_DIRECTORY"

# Cleanup everything
cleanup
