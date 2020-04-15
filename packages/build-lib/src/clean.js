#!/usr/bin/env node

const rimraf = require('rimraf');

rimraf('lib', () => {});
rimraf('./node_modules/.cache', () => {});
