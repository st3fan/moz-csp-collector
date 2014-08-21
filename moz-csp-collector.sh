#!/bin/sh

# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/

if [ -z "$DB_URL" ]; then
   DB_URL="postgres://csp:csp@localhost/csp"
fi

exec /usr/local/bin/moz-csp-collector -database="$DB_URL"
