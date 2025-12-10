#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

// eslint-disable-next-line import/no-extraneous-dependencies
const ts = require('typescript');

const DTS_FILE = path.resolve(__dirname, '..', 'node_modules', '@material-design-icons', 'font', 'index.d.ts');
const TYPE_NAME = 'MaterialIcons';
const OUTPUT_FILE = path.resolve(__dirname, '..', 'src', 'assets', 'material-icons', 'MaterialIcons.json');

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
  console.error(`❌ Failed to find tuple type "${TYPE_NAME}" in file ${DTS_FILE}`);
  process.exit(1);
}

const json = JSON.stringify(values, null, 2);

fs.writeFileSync(OUTPUT_FILE, json);

// eslint-disable-next-line no-console
console.log(`✅ JSON successfully generated: ${OUTPUT_FILE}`);
