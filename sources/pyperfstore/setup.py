import os
from setuptools import setup

# Utility function to read the README file.
# Used for the long_description.  It's nice, because now 1) we have a top level
# README file and 2) it's easier to type in the README file than to put a raw
# string in below ...
def read(fname):
    return open(os.path.join(os.path.dirname(__file__), fname)).read()

install_requires = [
	"pymongo",
]

setup(
	name = "pyperfstore",
	version = "0.0.1",
	author = "Capensis",
	author_email = "canopsis@capensis.fr",
	description = ("Store performance data"),
	license = "AGPL v3",
	zip_safe = False,
	keywords = "nagios perfdata storage store performance canopsis",
	install_requires=install_requires,
	url = "http://www.canopsis.org",
	packages=['pyperfstore'],
	scripts=['scripts/pyperfstore'],
	long_description=read('README'),
)
