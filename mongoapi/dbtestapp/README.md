<!--
Copyright 2021-present Open Networking Foundation
SPDX-License-Identifier: Apache-2.0
-->

# testapp

## test application used to demonstrate mongoDB library capabilities. 

### Capabilities shown
        - Inserting documents with timeout, so the document is removed from the collection after a certain period of time. 

## How to use ?

- clone aether-in-a-box repository
- cd aether-in-a-box
- make dbtestapp
- Now you should see dbtestapp running. You can use REST APIs to trigger dbtestapp code.
- Example of how to use REST APIs is provided in api-script.sh
