# -*- mode: python -*-

import os

from subprocess import check_output

from PyInstaller.utils.hooks import collect_data_files

raw_imports = check_output("./find_imports.sh")
imports = ["kombu.transport.pyamqp"]
imports.extend(raw_imports.split('\n'))
datas = []
datas.extend(collect_data_files('jsonschema'))

app_entry_script='../scripts/webserverpy'
app_bin_name='webserver'
app_dir_name='webserver-dir'

block_cipher = None
a = Analysis(
  [app_entry_script],
  pathex=['.'],
  hiddenimports=imports,
  hookspath=None,
  runtime_hooks=None,
  cipher=block_cipher,
  datas=datas
)

pyz = PYZ(a.pure, cipher=block_cipher)

exe = EXE(
  pyz,
  a.scripts,
  a.binaries,
  a.zipfiles,
  name=app_bin_name,
  debug=False,
  strip=True,
  upx=True,
  console=True
)
coll = COLLECT(
  exe,
  a.binaries,
  a.zipfiles,
  a.datas,
  strip=True,
  upx=True,
  name=app_dir_name
)
