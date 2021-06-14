#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

const loadEnv = require('./load-env'); // eslint-disable-line import/no-extraneous-dependencies

const localEnvPath = path.resolve(process.cwd(), '.env.local');
const baseEnvPath = path.resolve(process.cwd(), '.env');

loadEnv(localEnvPath);
loadEnv(baseEnvPath);

if (process.env.PREPARE_COMMIT_MSG_HOOK !== 'enabled') {
  process.exit();
}

function getPathToGitHead(folder) {
  return path.resolve(folder, '.git', 'HEAD');
}

function getPathToParentFolder(folder) {
  return path.resolve(folder, '..');
}

function isRepositoryRoot(folder) {
  return fs.existsSync(getPathToGitHead(folder));
}

function findRepositoryRoot() {
  let repositoryRoot = __dirname;

  while (!isRepositoryRoot(repositoryRoot) && fs.existsSync(repositoryRoot)) {
    repositoryRoot = getPathToParentFolder(repositoryRoot);
  }

  return repositoryRoot;
}

function getBranchName(repositoryRoot) {
  const head = fs.readFileSync(getPathToGitHead(repositoryRoot)).toString();
  const [, branchName] = head.match(/^ref: refs\/heads\/(.*)/) || [];

  return branchName;
}

const repositoryRoot = findRepositoryRoot();

if (!fs.existsSync(repositoryRoot)) {
  console.error('The script was unable to find the root of the Git repository.');
  process.exit(1);
}

const branchName = getBranchName(repositoryRoot);

if (!branchName) {
  process.exit();
}

const [, branchPrefix, issueNumber] = branchName.match(/^(.+)\/(#\d+)/i) || [];

if (!branchPrefix || !issueNumber) {
  process.exit();
}

const huskyGitParams = process.env.HUSKY_GIT_PARAMS;

if (!huskyGitParams) {
  console.error('The script expects Git parameters to be accessible via HUSKY_GIT_PARAMS.');
  process.exit(1);
}


const [commitMessageFile] = huskyGitParams.split(' ');

if (!commitMessageFile) {
  console.error('The script requires HUSKY_GIT_PARAMS to contain the name of the file containing the commit log message.');
  process.exit(1);
}

const commitPrefix = `${branchPrefix}(${issueNumber}): `;

const pathToCommitMessageFile = path.resolve(repositoryRoot, commitMessageFile);
const content = fs.readFileSync(pathToCommitMessageFile);

if (content.indexOf(commitPrefix) === -1) {
  fs.writeFileSync(pathToCommitMessageFile, commitPrefix + content);
}
