# -*- mode: python -*-

import os

from subprocess import check_output

from PyInstaller.utils.hooks import collect_data_files

"""
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
"""

app_entry_script=os.environ['PYI_SCRIPT']
app_bin_name=os.environ['PYI_BIN_NAME']

app_debug=False
if int(os.environ.get('PYI_DEBUG', 0)) == 1:
    app_debug=True

app_strip=False
if int(os.environ.get('PYI_STRIP', 0)) == 1:
    app_strip=True

app_upx=False
if int(os.environ.get('PYI_UPX', 0)) == 1:
    app_upx=True

a = Analysis(
    [app_entry_script],
    pathex=['./'],
    hookspath=['./hooks'],
    hiddenimports=[],
    runtime_hooks=[],
    excludes=[],
)

pyz = PYZ(a.pure, a.zipped_data)

exe = EXE(pyz, a.scripts, a.binaries, a.zipfiles,
    exclude_binaries=True,
    name=app_bin_name,
    debug=app_debug,
    strip=app_strip,
    upx=app_upx,
    console=True
)
coll = COLLECT(exe, a.binaries, a.zipfiles, a.datas,
    strip=app_strip,
    upx=app_upx,
    name=app_bin_name
)
