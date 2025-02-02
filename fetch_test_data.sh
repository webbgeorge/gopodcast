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
  if [[ $2 == "null" ]]; then
    echo "no feed url, skipping"
    return
  fi
  curl -sL -o "testdata/top-podcasts/${1}.xml" $2
}

mkdir -p testdata/top-podcasts

# use the top 50 UK, top 50 US, top 10 German and top 10 Japanese podcasts for tests
processChartURL "https://rss.marketingtools.apple.com/api/v2/gb/podcasts/top/50/podcasts.json"
processChartURL "https://rss.marketingtools.apple.com/api/v2/us/podcasts/top/50/podcasts.json"
processChartURL "https://rss.marketingtools.apple.com/api/v2/de/podcasts/top/10/podcasts.json"
processChartURL "https://rss.marketingtools.apple.com/api/v2/jp/podcasts/top/10/podcasts.json"
