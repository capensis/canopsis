# -*- mode: python -*-

from PyInstaller.utils.hooks import collect_data_files

block_cipher = None
datas = []
datas.extend(collect_data_files('jsonschema'))

a = Analysis(['scripts/engine-launcher'],
            pathex=['.'],
            hiddenimports=[
                "canopsis",
                "canopsis.engines.dynamic",
                "canopsis.stats.process",
                "canopsis.configuration.driver.file",
                "canopsis.configuration.driver.file.json",
                "canopsis.configuration.driver.file.ini",
                "kombu.transport.pyamqp",
            ],
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

