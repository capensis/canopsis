# -*- mode: python -*-

import os

from subprocess import check_output

from PyInstaller.utils.hooks import collect_data_files

app_entry_script=os.environ['PYI_SCRIPT']
app_bin_name=os.environ['PYI_BIN_NAME']
app_dir_name=os.environ['PYI_DIR_NAME']

a = Analysis(
  [app_entry_script],
  pathex=['.'],
  hookspath=['./pyinstaller/hooks'],
)

pyz = PYZ(a.pure)

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
