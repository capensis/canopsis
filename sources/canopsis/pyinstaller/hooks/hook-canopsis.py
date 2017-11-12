import os

from PyInstaller.utils.hooks import collect_data_files


def _find_canopsis_pyfiles(srcdir):
    canopsis_files = []
    for dirpath, subdirs, files in os.walk(srcdir):
        for file_ in [f for f in files if f.endswith('.py')]:
            canopsis_files.append(os.path.join(dirpath, file_))

    return canopsis_files


def _find_canopsis_imports(srcdir):
    pyfiles = _find_canopsis_pyfiles(srcdir)
    imports = []
    for pyfile in pyfiles:
        if pyfile.endswith('__init__.py'):
            continue

        import_ = pyfile.replace(os.sep, '.').replace('.py', '')

        if '.cli.' in import_:
            continue

        imports.append(import_)

    return imports


def get_static_hidden_imports():
    imports = [
        'gunicorn.workers.ggevent',
        'gunicorn.glogging',
        'kombu.transport.pyamqp',
        'validictory'
    ]

    return imports


def get_additional_data():
    return collect_data_files('jsonschema')


#datas = get_additional_data()
#hiddenimports = get_static_hidden_imports()
hiddenimports = _find_canopsis_imports('canopsis')
