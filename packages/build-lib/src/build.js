#!/usr/bin/env node

const rollup = require("rollup");
const rollupTypescript = require("rollup-plugin-typescript2");
const __typescript = require("typescript");
const sass = require('rollup-plugin-sass')

const commonjs = require("@rollup/plugin-commonjs");
const resolve = require("@rollup/plugin-node-resolve");
const excludeDependenciesFromBundle = require("rollup-plugin-exclude-dependencies-from-bundle");

// see below for details on the options
const inputOptions = {
  input: "src/index.ts",
  plugins: [
    excludeDependenciesFromBundle(),
    rollupTypescript({
      tsconfig: "tsconfig.json",
      typescript: __typescript,
    }),
    sass({
      output: true,
    }),
    resolve(),
    commonjs(),
  ],
};

const outputOptions = [
  {
    file: "./lib/index.js",
    format: "cjs",
  },
  {
    file: "./lib/index.es.js",
    format: "es",
  },
];

const generateBundles = (bundle) =>
  Promise.all(
    outputOptions.map((outputOption) => {
      console.log(`Writing ${outputOption.file}`);
      return bundle.write(outputOption);
    })
  ).then(() => bundle);

async function build() {
  // create a bundle
  const bundle = await rollup
    .rollup(inputOptions)
    .then(generateBundles)
    .catch((error) => {
      console.error(error);
      return process.exit(-1);
    });
}

build();
