#!/bin/bash

echo running vulnerability scanner
govulncheck --format openvex ./... | jq '.statements[] | select (.status == "affected").vulnerability'

echo "NOTE: The vulnerability output above is informational. Please review the output to see if it is possible to patch the vulnerability."