import { spawnSync } from 'node:child_process';
import { readdirSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const projectRoot = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const patchesDirectory = resolve(projectRoot, 'patches');
const patches = readdirSync(patchesDirectory)
  .filter(file => file.endsWith('.patch'))
  .sort();

const runPatch = (patch, ...arguments_) =>
  spawnSync('patch', [...arguments_, '-p1', '-i', patch], {
    cwd: projectRoot,
    encoding: 'utf8',
  });

for (const patch of patches) {
  const patchPath = resolve(patchesDirectory, patch);
  const check = runPatch(patchPath, '--dry-run', '--batch', '--forward');

  if (check.status === 0) {
    const apply = runPatch(patchPath, '--batch', '--forward');

    if (apply.status !== 0) {
      throw new Error(`Unable to apply ${patch}:\n${apply.stderr}`);
    }

    console.log(`${patch} applied`);
    continue;
  }

  const alreadyApplied = runPatch(patchPath, '--dry-run', '--batch', '--reverse').status === 0;

  if (!alreadyApplied) {
    throw new Error(`Unable to apply ${patch}:\n${check.stderr}`);
  }

  console.log(`${patch} already applied`);
}
