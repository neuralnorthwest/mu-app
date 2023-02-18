#!/bin/bash

# Copyright 2023 Scott M. Long
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# 	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Check for files missing license headers and not in .licenseignore.

set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

NO_LIC_FILES=$(grep --exclude-dir=.git --exclude-dir=venv --exclude-dir=gen --exclude-dir=mock --exclude-dir=docs -HLr 'Licensed under the Apache License' . | sort)
IGNORE_FILES=$(cat .licenseignore | sort)
REPORT_FILES=""
while read -r NO_LIC_FILE; do
    if ! grep -q "$NO_LIC_FILE" <<< "$IGNORE_FILES"; then
        REPORT_FILES="$REPORT_FILES $NO_LIC_FILE"
    fi
done <<< "$NO_LIC_FILES"

if [ -n "$REPORT_FILES" ]; then
    echo "Files missing license headers:"
    echo "$REPORT_FILES"
    exit 1
fi
