#!/usr/bin/env node

import { execFileSync } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';

const HOOK_ENV_VARIABLE = 'PREPARE_COMMIT_MSG_HOOK';

function parseEnvValue(content, variable) {
  const line = content
    .split(/\r?\n/u)
    .find(value => new RegExp(`^\\s*(?:export\\s+)?${variable}\\s*=`, 'u').test(value));

  if (!line) {
    return undefined;
  }

  const value = line.slice(line.indexOf('=') + 1).trim();
  const quote = value[0];

  if ((quote === '"' || quote === '\'') && value.endsWith(quote)) {
    return value.slice(1, -1);
  }

  return value.replace(/\s+#.*$/u, '').trim();
}

function getEnvValue(variable) {
  if (process.env[variable] !== undefined) {
    return process.env[variable];
  }

  for (const filename of ['.env.local', '.env']) {
    const envPath = path.resolve(process.cwd(), filename);

    if (!fs.existsSync(envPath)) {
      continue;
    }

    const value = parseEnvValue(fs.readFileSync(envPath, 'utf8'), variable);

    if (value !== undefined) {
      return value;
    }
  }

  return undefined;
}

if (getEnvValue(HOOK_ENV_VARIABLE) !== 'enabled') {
  process.exit();
}

const [, , commitMessageFile] = process.argv;

if (!commitMessageFile) {
  console.error('The script requires the path to the commit message file.');
  process.exit(1);
}

const repositoryRoot = execFileSync('git', ['rev-parse', '--show-toplevel'], {
  encoding: 'utf8',
}).trim();
const branchName = execFileSync('git', ['branch', '--show-current'], {
  encoding: 'utf8',
}).trim();
const [, branchPrefix, issueNumber] = branchName.match(/^(.+)\/(#\d+)/iu) || [];

if (!branchPrefix || !issueNumber) {
  process.exit();
}

const commitPrefix = `${branchPrefix}(${issueNumber}): `;
const commitMessagePath = path.resolve(repositoryRoot, commitMessageFile);
const content = fs.readFileSync(commitMessagePath, 'utf8');

if (!content.includes(commitPrefix)) {
  fs.writeFileSync(commitMessagePath, commitPrefix + content);
}
