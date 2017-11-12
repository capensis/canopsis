# -*- mode: python -*-

import os

from subprocess import check_output

from PyInstaller.utils.hooks import collect_data_files

def import_hook():
    try:
        import importlib.util
        spec = importlib.util.spec_from_file_location('hook_canopsis', 'hooks/hook-canopsis.py')
        hook_canopsis = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(foo)
    except ImportError:
        import imp
        hook_canopsis = imp.load_source('hook_canopsis', 'hooks/hook-canopsis.py')

    return hook_canopsis

hook_canopsis = import_hook()

app_entry_script=os.environ['PYI_SCRIPT']
app_bin_name=os.environ['PYI_BIN_NAME']
app_dir_name=os.environ['PYI_DIR_NAME']
app_strip=False
if int(os.environ.get('PYI_STRIP', 0)) == 1:
    app_strip=True

a = Analysis(
    [app_entry_script],
    pathex=['.'],
    hookspath=['./hooks'],
    datas=hook_canopsis.get_additional_data(),
    #hiddenimports=hook_canopsis.get_static_hidden_imports(),
)

pyz = PYZ(a.pure)

exe = EXE(
    pyz,
    a.scripts,
    a.binaries,
    a.zipfiles,
    name=app_bin_name,
    debug=False,
    strip=app_strip,
    upx=True,
    console=True,
    append_pkg=False
)
coll = COLLECT(
    exe,
    a.binaries,
    a.zipfiles,
    a.datas,
    strip=app_strip,
    upx=True,
    name=app_dir_name
)
