#!/bin/bash
#
# Fetch feeds of most popular podcasts to be used as test data
#
# Dependencies: curl and jq

set -e

function processChartURL() {
  curl -sL $1 | jq -r .feed.results[].id | while read ID; do
    echo "Podcast ${ID}"
    RES=$(curl -sL "https://itunes.apple.com/lookup?id=${ID}")
    echo $RES | jq -r .results[0].collectionName
    downloadTestDataFeed $ID $(echo $RES | jq -r .results[0].feedUrl)
  done
}

function downloadTestDataFeed() {
  curl -sL -o "testdata/top-podcasts/${1}.xml" $2
}

mkdir -p testdata/top-podcasts

processChartURL "https://rss.marketingtools.apple.com/api/v2/gb/podcasts/top/50/podcasts.json"
