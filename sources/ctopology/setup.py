from setuptools import setup

import ctopology

install_requires = [
	"cstorage",
]

with open('README') as f:
	desc = f.read()

setup(
	name=ctopology.__name__,
	version=ctopology.__version__,
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("Store topology"),
	license="AGPL v3",
	keywords="ctopology storage store canopsis ctopology",
	install_requires=install_requires,
	url="http://www.canopsis.org",
	packages=['ctopology'],
	scripts=['scripts/ctopology'],
	long_description=desc,
	test_suite="test"
)
