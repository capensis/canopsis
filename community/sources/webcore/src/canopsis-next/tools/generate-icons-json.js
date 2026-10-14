#!/usr/bin/env node

import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

import ts from 'typescript';

const dirname = path.dirname(fileURLToPath(import.meta.url));

const DTS_FILE = path.resolve(dirname, '..', 'node_modules', 'material-symbols', 'index.d.ts');
const TYPE_NAME = 'MaterialSymbols';
const OUTPUT_FILE = path.resolve(dirname, '..', 'src', 'assets', 'material-icons', 'MaterialIcons.json');

const content = fs.readFileSync(DTS_FILE, 'utf8');

const source = ts.createSourceFile(
  DTS_FILE,
  content,
  ts.ScriptTarget.Latest,
  true,
);

let values = null;

ts.forEachChild(source, (node) => {
  if (
    ts.isTypeAliasDeclaration(node)
    && node.name.text === TYPE_NAME
  ) {
    const typeNode = node.type;

    if (ts.isTupleTypeNode(typeNode)) {
      values = typeNode.elements.map(el => el.literal.text);
    }
  }
});

if (!values) {
  console.error(`Failed to find tuple type "${TYPE_NAME}" in file ${DTS_FILE}`);
  process.exit(1);
}

const json = JSON.stringify(values, null, 2);

fs.writeFileSync(OUTPUT_FILE, json);

console.log(`JSON successfully generated: ${OUTPUT_FILE}`);
