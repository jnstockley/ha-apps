#!/bin/sh

set -eu

config_file="doh/config.yaml"
current_version="$(sed -nE 's/^version:[[:space:]]*"([0-9]+\.[0-9]+\.[0-9]+)".*/\1/p' "${config_file}")"

if [ -z "${current_version}" ]; then
    echo "Could not find a semantic version in ${config_file}" >&2
    exit 1
fi

major_version="${current_version%%.*}"
remaining_version="${current_version#*.}"
minor_version="${remaining_version%%.*}"
patch_version="${remaining_version#*.}"
next_version="${major_version}.${minor_version}.$((patch_version + 1))"
temp_file="$(mktemp "${config_file}.XXXXXX")"
trap 'rm -f "${temp_file}"' EXIT

awk -v version="${next_version}" '
  $1 == "version:" {
    print "version: \"" version "\""
    updated = 1
    next
  }
  { print }
  END {
    if (!updated) {
      exit 1
    }
  }
' "${config_file}" > "${temp_file}"

mv "${temp_file}" "${config_file}"
trap - EXIT

echo "Bumped DoH app version from ${current_version} to ${next_version}"
