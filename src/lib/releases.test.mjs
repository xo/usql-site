// Tests of src/lib/releases.mjs. Run them with:
//
//   node --test "src/lib/*.test.mjs"

import assert from "node:assert/strict";
import { test } from "node:test";
import { notes, safeURL } from "./releases.mjs";

/** walk returns every node of the trees of nodes. */
function walk(nodes) {
  return nodes.flatMap((n) => [n, ...walk(n.children ?? [])]);
}

/** content returns the text of the trees of nodes. */
function content(nodes) {
  return walk(nodes)
    .map((n) => n.value ?? "")
    .join("");
}

const hostile = `# Top

<script>alert(1)</script>

Inline <img src="https://evil.example/x.png"> and <b>bold</b>.

![remote](https://evil.example/a.png) ![local](/images/x.png) ![](https://evil.example/noalt.png)

[![badge](https://img.shields.io/x.svg)](https://github.com/xo/usql)

![ref image][img] and [ref link][lnk] and [bad][bad]

[img]: https://evil.example/ref.png
[lnk]: https://example.com/ok
[bad]: javascript:alert(2)

[js](javascript:alert(1)) [rel](docs/x.md) [frag](#top) <https://auto.example/>

::cast[x]{name="tour"}

- <iframe src="https://evil.example/"></iframe>

| a |
|---|
| ![t](https://evil.example/t.png) |
`;

test("notes leave no raw HTML, image or reference", () => {
  const nodes = walk(notes(hostile));
  for (const type of [
    "html",
    "image",
    "imageReference",
    "linkReference",
    "definition",
    "leafDirective",
  ]) {
    assert.deepEqual(
      nodes.filter((n) => n.type === type),
      [],
      `a node of type ${type} is left`
    );
  }
});

test("notes keep raw HTML as text", () => {
  const text = content(notes(hostile));
  assert.match(text, /<script>alert\(1\)<\/script>/);
  assert.match(text, /<iframe src="https:\/\/evil.example\/"><\/iframe>/);
});

test("notes keep only links to absolute http, https and mailto URLs", () => {
  const links = walk(notes(hostile)).filter((n) => n.type === "link");
  assert.deepEqual(
    links.map((n) => n.url),
    [
      "https://evil.example/a.png",
      "https://evil.example/noalt.png",
      "https://github.com/xo/usql",
      "https://evil.example/ref.png",
      "https://example.com/ok",
      "https://auto.example/",
      "https://evil.example/t.png",
    ]
  );
  const badge = links.find((n) => n.url === "https://github.com/xo/usql");
  assert.deepEqual(badge?.children, [{ type: "text", value: "badge" }]);
});

test("notes move headings to level 4 or lower", () => {
  const depths = walk(notes("# One\n\n## Two\n\n### Three\n\n###### Six"))
    .filter((n) => n.type === "heading")
    .map((n) => n.depth);
  assert.deepEqual(depths, [4, 4, 5, 6]);
});

test("safeURL", () => {
  assert.equal(safeURL("https://github.com/"), true);
  assert.equal(safeURL("mailto:a@example.com"), true);
  assert.equal(safeURL("javascript:alert(1)"), false);
  assert.equal(safeURL("/images/x.png"), false);
  assert.equal(safeURL(""), false);
});
