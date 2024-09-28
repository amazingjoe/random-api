/*
Copyright (C) 2024  Kaan Barmore-Genc

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
import { hydrate, prerender as ssr } from "preact-iso";
import { ComponentChildren } from "preact";

import "./style.css";
import { Signal, signal } from "@preact/signals";
import { useCallback } from "preact/hooks";

//#region Utils

const colors = [
  "text-red-500",
  "text-sky-500",
  "text-yellow-500",
  "text-green-500",
  "text-blue-500",
  "text-indigo-500",
  "text-rose-500",
  "text-purple-500",
  "text-pink-500",
  "text-cyan-500",
  "text-emerald-500",
];
const keywordColors = new Map<string, string>();
function keywordColor(keyword: string) {
  let color = keywordColors.get(keyword);
  if (!color) {
    color = colors.pop();
    if (!color) {
      throw new Error("Out of colors");
    }
    keywordColors.set(keyword, color);
  }
  return color;
}

function clsx(...classes: (string | boolean | undefined)[]) {
  return classes.filter(Boolean).join(" ");
}

function C({
  children: keyword,
  pad = false,
}: {
  children: string;
  pad?: boolean;
}) {
  return (
    <code class={clsx(pad && "p-1", keywordColor(keyword))}>{keyword}</code>
  );
}
//#endregion

//#region Components
function Header() {
  return (
    <header class="w-full bg-black bg-opacity-20 px-12 py-8 flex flex-row items-center">
      <h2>Random Generation API</h2>
      <div class="grow"></div>
      <a
        class="flex flex-row items-center gap-4 btn btn-outline"
        href="https://github.com/SeriousBug/random-api"
      >
        Github
        {/* Icon from https://heroicons.com, MIT license */}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth="1.5"
          stroke="currentColor"
          class="size-6"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25"
          />
        </svg>
      </a>
    </header>
  );
}

function Footer() {
  return (
    <footer class="w-full bg-black bg-opacity-20 px-12 py-8 flex flex-row justify-center">
      <h2 class="flex flex-row">
        Made with{" "}
        <span class="inline-block heart bg-red-600" aria-label="love"></span> in
        Illinois
      </h2>
      <div class="bg-base-content mx-4 circle"></div>
      <a href="https://github.com/SeriousBug/random-api">
        This is open source software, licensed under AGPLv3.
      </a>
    </footer>
  );
}

function Layout({ children = null }: { children?: ComponentChildren }) {
  return (
    <div class="min-h-lvh flex flex-col">
      <Header />
      <main class="flex flex-col">{children}</main>
      <Footer />
    </div>
  );
}

type ParameterBase = {
  name: string;
};
type ParameterText = {
  type: "text";
  value?: string;
} & ParameterBase;
type ParameterNumber = {
  type: "number";
  step?: string;
  value?: string;
} & ParameterBase;
type ParameterEnum = {
  type: "enum";
  options: string[];
  value?: string;
} & ParameterBase;
type Parameter = ParameterText | ParameterNumber | ParameterEnum;

function parameterId({ name }: Pick<Parameter, "name">) {
  return `param-${name}`;
}

function ParameterDisplay({ param }: { param: Signal<Parameter> }) {
  const { name, value } = param.value;

  return (
    <label for={parameterId({ name })} key={name} class={value ? "" : "hidden"}>
      <C pad={false}>{name}</C>
      <span class="text-opacity-80 text-base-content">=</span>
      <span class="text-base-content">{value}</span>
    </label>
  );
}

function ParameterInputInner({ param }: { param: Signal<Parameter> }) {
  const { type } = param.value;

  if (type === "text") {
    return (
      <input
        onInput={(e) => {
          param.value = { ...param.value, value: e.currentTarget.value };
        }}
        type="text"
        id={parameterId(param.value)}
        class="input input-primary"
      />
    );
  }

  if (type === "number") {
    return (
      <input
        onInput={(e) => {
          param.value = { ...param.value, value: e.currentTarget.value };
        }}
        type="number"
        step={param.value.step}
        id={parameterId(param.value)}
        class="input input-primary"
      />
    );
  }

  if (type === "enum") {
    return (
      <select
        id={parameterId(param.value)}
        class="select select-primary"
        onInput={(e) => {
          param.value = { ...param.value, value: e.currentTarget.value };
        }}
      >
        <option value=""></option>
        {param.value.options.map((value) => (
          <option key={value} value={value}>
            {value}
          </option>
        ))}
      </select>
    );
  }

  throw new Error(`Unknown parameter type: ${type}`);
}

function ParameterLabel({ param }: { param: Signal<Parameter> }) {
  return (
    <label for={parameterId(param.value)}>
      <C>{param.value.name}</C>
    </label>
  );
}

function ParameterInput({ param }: { param: Signal<Parameter> }) {
  return (
    <div class="w-full flex flex-col gap-2 my-2">
      <ParameterLabel param={param} />
      <ParameterInputInner param={param} />
    </div>
  );
}

function QueryStringMarker({
  parameters,
}: {
  parameters: Signal<Parameter>[];
}) {
  const hasParameterSelected = parameters.some((param) => !!param.value.value);

  if (!hasParameterSelected) {
    return null;
  }

  return <span class="text-opacity-80">?</span>;
}

const API_BASE = import.meta.env.API_URL ?? `https://rnd.bgenc.dev`;

function makeUrl({
  path,
  parameters,
}: {
  path: string;
  parameters: Signal<Parameter>[];
}) {
  const query = new URLSearchParams();
  parameters.forEach((param) => {
    if (param.value.value) {
      query.append(param.value.name, param.value.value);
    }
  });
  return `${API_BASE}${path}?${query.toString()}`;
}

function GenerateButton({
  path,
  method,
  parameters,
  result,
}: {
  path: string;
  method: string;
  parameters: Signal<Parameter>[];
  result: Signal<string>;
}) {
  const submit = useCallback(async () => {
    const out = await fetch(makeUrl({ path, parameters }), {
      method,
    });
    result.value = await out.text();
  }, []);

  return (
    <button onClick={submit} class="btn btn-primary w-48 self-center">
      Generate
    </button>
  );
}

function Output({ result }: { result: Signal<string> }) {
  return (
    <div class="w-full flex flex-col gap-2">
      <h2 class="text-lg">Output</h2>
      <pre class="p-4 bg-white bg-opacity-10 rounded items-center w-full text-center break-all text-wrap">
        {result.value ? result.value : " "}
      </pre>
    </div>
  );
}

function InfoModal({
  children,
  name,
}: {
  children: ComponentChildren;
  name: string;
}) {
  const id = `modal-${name}`;

  const onOpen = useCallback(() => {
    (document.getElementById(id) as HTMLDialogElement).showModal();
  }, [id]);
  const onClose = useCallback(() => {
    (document.getElementById(id) as HTMLDialogElement).close();
  }, [id]);

  return (
    <>
      <button className="btn btn-outline w-48 self-center" onClick={onOpen}>
        Learn more
      </button>
      <dialog
        id={id}
        className="modal"
        onMouseDown={(e) => {
          if ("nodeName" in e.target && e.target.nodeName === "DIALOG") {
            (e.target as HTMLDialogElement).close();
          }
        }}
      >
        <div className="modal-box flex flex-col gap-4">
          {children}
          <button className="btn btn-outline" onClick={onClose}>
            Close
          </button>
        </div>
      </dialog>
    </>
  );
}

function Endpoint({
  name,
  className,
  subtitle,
  method = "GET",
  path,
  parameters = [],
  children,
}: {
  name: string;
  className?: string;
  subtitle?: string;
  method?: string;
  path: string;
  parameters: Signal<Parameter>[];
  children?: ComponentChildren;
}) {
  const result = signal<string>("");
  const copyUrl = useCallback(() => {
    navigator.clipboard.writeText(makeUrl({ path, parameters }));
  }, [path, parameters]);

  return (
    <section
      class={clsx(
        "card bg-neutral text-neutral-content min-w-[480px] grow max-w-[540px] flex flex-col gap-4 p-8",
        className
      )}
    >
      <h2 class="text-2xl">{name}</h2>
      {subtitle ? <h3 class="opacity-90">{subtitle}</h3> : null}
      <div class="flex flex-row gap-2">
        <div class="text-opacity-80 p-4 rounded bg-white bg-opacity-10 font-mono">
          {method}
        </div>
        <button
          class="p-4 rounded bg-white bg-opacity-10 font-mono w-full break-all"
          onClick={copyUrl}
        >
          <span>
            {API_BASE}
            {path}
          </span>
          <span>
            <QueryStringMarker parameters={parameters} />
            {parameters.map((param, i) => (
              <ParameterDisplay key={i} param={param} />
            ))}
          </span>
        </button>
      </div>
      {parameters.map((param, i) => (
        <ParameterInput key={i} param={param} />
      ))}
      <Output result={result} />
      <GenerateButton
        path={path}
        method={method}
        parameters={parameters}
        result={result}
      />
      {children ? <InfoModal name={name}>{children}</InfoModal> : null}
    </section>
  );
}
//#endregion

function Endpoints() {
  return (
    <div class="flex flex-rows flex-wrap gap-8 grow items-start justify-center p-8 max-w-[1800px] mx-auto">
      <Endpoint
        name="Integer"
        subtitle="Returns a random integer, a subset of whole numbers."
        path="/v1/int"
        parameters={[
          signal<Parameter>({ type: "number", name: "min" }),
          signal<Parameter>({ type: "number", name: "max" }),
        ]}
      >
        <p>
          The limits are inclusive of <C>min</C>, but exclusive of <C>max</C>.
          So for example, with <C>min</C> set to 0 and <C>max</C> set to 10, the
          output will be between 0 and 9.
        </p>
        <p>
          If no <C>min</C> is specified, the default is 0. If no <C>max</C> is
          specified, the default is 100. You may use negative numbers for{" "}
          <C>min</C> and <C>max</C>. <C>min</C> must be strictly less than{" "}
          <C>max</C>.
        </p>
        <p>
          These limits are also bound by minimum and maximum values of go&apos;s
          int type. This practically means 32-bit integers.
        </p>
      </Endpoint>
      <Endpoint
        name="Floating-Point Number"
        subtitle="Returns a random floating-point number, a subset of real numbers."
        path="/v1/float"
        parameters={[
          signal<Parameter>({ type: "number", name: "min", step: "any" }),
          signal<Parameter>({ type: "number", name: "max", step: "any" }),
        ]}
      >
        <p>Returns a floating point number, a subset of real numbers.</p>
        <p>
          If no <C>min</C> is specified, the default is 0. If no <C>max</C> is
          specified, the default is 1. You may use negative numbers for{" "}
          <C>min</C> and <C>max</C>. <C>min</C> must be strictly less than{" "}
          <C>max</C>.
        </p>
      </Endpoint>
      <Endpoint
        name="Dice"
        subtitle="Rolls dice in a variety of formats, like 2d8."
        path="/v1/dice"
        parameters={[
          signal<Parameter>({ type: "text", name: "input" }),
          signal<Parameter>({
            type: "enum",
            name: "output",
            options: ["sum", "full"],
          }),
        ]}
      >
        <p>
          Generates a random number following an RPG dice format. Use the{" "}
          <C>input</C> parameter to specify what dice you are trying to roll. If
          not specified, the default is 1d6.
        </p>
        <ul class="list-disc pl-6">
          <li>
            Standard:{" "}
            <span class="bg-white bg-opacity-20 p-1 bordered font-mono">
              xdy[[k|d][h|l]z][+/-c]
            </span>{" "}
            rolls and sums x y-sided dice, keeping or dropping the lowest or
            highest z dice and optionally adding or subtracting c. Example:
            4d6kh3+4 means roll 4 six-sided dice, keep the highest 3, and add 4.
          </li>
          <li>
            Fudge:{" "}
            <span class="bg-white bg-opacity-20 p-1 bordered font-mono">
              xdf[+/-c]
            </span>{" "}
            rolls and sums x fudge dice (Dice that returns numbers between -1
            and 1), and optionally adding or subtracting c. Example: 4df+4
          </li>
          <li>
            Versus:{" "}
            <span class="bg-white bg-opacity-20 p-1 bordered font-mono">
              xdy[e|r]vt
            </span>{" "}
            rolls x y-sided dice, counting the number that roll t or greater.
          </li>
          <li>
            EotE:{" "}
            <span class="bg-white bg-opacity-20 p-1 bordered font-mono">
              xc [xc ...]
            </span>{" "}
            rolls x dice of color c (b, blk, g, p, r, w, y) and returns the
            aggregate result.
          </li>
          <li>
            Adding an e to the versus rolls above makes dice explode. Dice are
            rerolled and have the rolled value added to their total when they
            roll a y. Adding an r makes dice rolling a y add another die to the
            pool instead.
          </li>
        </ul>
        <p>
          Use the <C>output</C> parameter to specify what you want to get back.
          If not specified, the default is sum, which adds up the totals of all
          the dice rolled. If you specify full, you will get back not just the
          sum but a list of individual rolls, including discarded ones if there
          are any.
        </p>
      </Endpoint>
      <Endpoint name="ULID" path="/v1/ulid" parameters={[]}>
        <p>
          Generate a{" "}
          <a href="https://github.com/ulid/spec">
            Universally Unique Lexicographically Sortable Identifier
          </a>
          , a more compact and sortable alternative to UUID. There are no
          parameters, the output is always a 26 character string.
        </p>
      </Endpoint>
      <Endpoint
        name="UUID"
        path="/v1/uuid"
        parameters={[
          signal<Parameter>({
            type: "enum",
            name: "version",
            options: ["4", "7"],
          }),
        ]}
      >
        <p>
          Generate a{" "}
          <a href="https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)">
            Universally Unique Identifier
          </a>
          . Only <C>version</C>s 4 and 7 are supported. 4 is completely random,
          while 7 starts with a timestamp that allows it to be sortable by time.
          If not specified, the default <C>version</C> is 4.
        </p>
      </Endpoint>
      <Endpoint
        name="Nano ID"
        path="/v1/nanoid"
        parameters={[signal<Parameter>({ type: "number", name: "size" })]}
      >
        <p>
          An even more compact random identifier than ULID, with a customizable{" "}
          <C>size</C>. The default <C>size</C> is 21 characters, and it must
          always be 1 or higher. Differently from ULID, the nanoids are not
          sortable.
        </p>
      </Endpoint>
      <Endpoint
        name="Word"
        subtitle="Generate one or more random words, from a variety of categories."
        path="/v1/word"
        className="max-w-[840px]"
        parameters={[
          signal<Parameter>({
            type: "enum",
            name: "category",
            options: [
              "words",
              "animals",
              "cities",
              "countries",
              "fruits",
              "vegetables",
              "lorem-ipsum",
              "nouns",
            ],
          }),
          signal<Parameter>({ type: "number", name: "count" }),
          signal<Parameter>({ type: "text", name: "separator" }),
        ]}
      >
        <p>
          Use <C>category</C> to pick what type of word you want.
        </p>
        <ul class="list-disc pl-6">
          <li>
            <span class="font-bold">words:</span> The default category. This is
            a list of over 9000 common English words, all at least 3 letters
            long. There should be no swear words in this list. There are proper
            nouns in this list.
          </li>
          <li>
            <span class="font-bold">
              animals: Almost 300 animal names in English.
            </span>
          </li>
          <li>
            <span class="font-bold">cities:</span> A list of over 2500 city
            names from around the world. This list comes from United Nations
            list of cities with 100,000 or more inhabitants. Out of that list,
            all city names that contain spaces or symbols were removed.
          </li>
          <li>
            <span class="font-bold">countries:</span> A list of 195 countries.
            This list may include some &quot;controversial&quot; countries that
            are not universally recognized, such as Taiwan, and omits others
            like Northern Cyprus. Some of these countries have spaces in their
            names, in which case they still count as one word despite the space.
          </li>
          <li>
            <span class="font-bold">fruits:</span> A list of around 50 fruits,
            only including fruits that have single word names that I recognize.
            This list is not exhaustive by any means. It only includes fruits by
            the culinary definition, ones like tomatoes are not included despite
            their botanical classification.
          </li>
          <li>
            <span class="font-bold">vegetables:</span> A list of around 50
            vegetables, only including vegetables that have single word names
            that I recognize. This list is not exhaustive by any means. It only
            includes vegetables by the culinary definition, ones like cucumbers
            are not included despite their botanical classification.
          </li>
          <li>
            <span class="font-bold">lorem-ipsum:</span> Around 140 words that
            are commonly used in &quot;lorem ipsum&quot; placeholder text.
          </li>
          <li>
            <span class="font-bold">nouns</span> A list of 1000 most common
            English nouns. Includes abbreviations like TV, but no proper nouns.
          </li>
        </ul>
        <p>
          The <C>count</C> parameter specifies how many words you want. This
          must be 1 or higher. The default is 1. If the <C>count</C> is set to
          higher than 1, the words will be separated by the <C>separator</C>, a
          space by default.
        </p>
      </Endpoint>
    </div>
  );
}

export function App() {
  return (
    <Layout>
      <Endpoints />
    </Layout>
  );
}

if (typeof window !== "undefined") {
  hydrate(<App />, document.getElementById("app"));
}

export async function prerender(data) {
  return await ssr(<App {...data} />);
}
