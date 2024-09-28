#! /usr/bin/env node
import fs from "fs";

const animalsRaw = fs.readFileSync("animals.raw.txt", "utf-8");

const animals = [
  ...new Set(
    animalsRaw
      .split("\n")
      .filter(
        (w) =>
          !w.match(/[0-9'\/_-]/) &&
          [...w.matchAll(" ")].length === 0 &&
          w.match(/[a-z]/),
      )
      .map((w) => w.replace(/[ ]?\(.*\)/, "")),
  ),
];

fs.writeFileSync("animals.txt", animals.join("\n"));

console.log(`Got ${animals.length} animals`);
