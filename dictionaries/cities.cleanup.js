#! /usr/bin/env node
import fs from "fs";

const citiesRaw = fs.readFileSync("cities.raw.txt", "utf-8");

const cities = [
  ...new Set(
    citiesRaw
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

fs.writeFileSync("cities.txt", cities.join("\n"));

console.log(`Got ${cities.length} cities`);
