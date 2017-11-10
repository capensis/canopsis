python setup.py build install
pyinstaller -y --clean -D pyinstaller/engine-launcher.spec
pyinstaller -y --clean -D pyinstaller/webserver.spec
