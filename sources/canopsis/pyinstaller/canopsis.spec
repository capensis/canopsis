# -*- mode: python -*-

import os

from subprocess import check_output

from PyInstaller.utils.hooks import collect_data_files

try:
    import importlib.util
    spec = importlib.util.spec_from_file_location('hook_canopsis', 'pyinstaller/hooks/hook-canopsis.py')
    hook_canopsis = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(foo)
except ImportError:
    import imp
    hook_canopsis = imp.load_source('hook_canopsis', 'pyinstaller/hooks/hook-canopsis.py')

app_entry_script=os.environ['PYI_SCRIPT']
app_bin_name=os.environ['PYI_BIN_NAME']
app_dir_name=os.environ['PYI_DIR_NAME']

a = Analysis(
    [app_entry_script],
    pathex=['.'],
    hookspath=['./pyinstaller/hooks'],
    datas=hook_canopsis.get_additional_data(),
    hiddenimports=hook_canopsis.get_static_hidden_imports(),
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
    console=True,
    append_pkg=False
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
