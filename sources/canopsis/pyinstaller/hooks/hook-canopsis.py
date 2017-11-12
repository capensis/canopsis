import os

from PyInstaller.utils.hooks import collect_data_files
from PyInstaller.compat import modname_tkinter


def _find_canopsis_pyfiles(srcdir):
    olddir = os.path.abspath(os.curdir)
    os.chdir(srcdir)
    canopsis_files = []
    for dirpath, subdirs, files in os.walk('.'):
        for file_ in [f for f in files if f.endswith('.py')]:
            canopsis_files.append(os.path.join(dirpath, file_))

    os.chdir(olddir)

    return canopsis_files


def _find_canopsis_imports(srcdir):
    pyfiles = _find_canopsis_pyfiles(srcdir)
    imports = []
    for pyfile in pyfiles:
        if pyfile.endswith('__init__.py'):
            continue

        import_ = pyfile.replace(os.sep, '.').replace('.py', '')
        import_ = import_.replace('..', 'canopsis.')

        if '.cli.' in import_:
            continue

        imports.append(import_)

    return imports


def get_static_hidden_imports():
    imports = [
        'gunicorn.workers.ggevent',
        'gunicorn.glogging',
        'kombu.transport.pyamqp',
        'validictory',
    ]

    return imports


def get_additional_data():
    return collect_data_files('jsonschema')


datas = get_additional_data()
hiddenimports = get_static_hidden_imports()
hiddenimports.extend(_find_canopsis_imports('../canopsis'))
excludedimports = [modname_tkinter]
