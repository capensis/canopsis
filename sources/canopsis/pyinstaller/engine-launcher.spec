# -*- mode: python -*-

import os

from subprocess import check_output

from PyInstaller.utils.hooks import collect_data_files

raw_imports = check_output("./find_imports.sh")
imports = ["kombu.transport.pyamqp"]
imports.extend(raw_imports.split('\n'))
datas = []
datas.extend(collect_data_files('jsonschema'))

block_cipher = None
a = Analysis(['../scripts/engine-launcher'],
            pathex=['.'],
            hiddenimports=imports,
            hookspath=None,
            runtime_hooks=None,
            cipher=block_cipher,
            datas=datas
)

pyz = PYZ(a.pure, cipher=block_cipher)

exe = EXE(pyz,
          a.scripts,
          a.binaries,
          a.zipfiles,
          name='engine-launcher',
          debug=False,
          strip=False,
          upx=False,
          console=True)

