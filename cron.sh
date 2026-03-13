#!/usr/bin/env sloth

with ($cwd = "${$scriptDir}/tool") {
  ["go", "run", "."]!
}
