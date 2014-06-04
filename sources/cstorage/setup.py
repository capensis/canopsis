from setuptools import setup

import cstorage

install_requires = [
	"pymongo",
]

with open('README') as f:
	desc = f.read()

setup(
	name=cstorage.__name__,
	version=cstorage.__version__,
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("Store data"),
	license="AGPL v3",
	keywords="storage store canopsis",
	install_requires=install_requires,
	url="http://www.canopsis.org",
	packages=[cstorage.__name__],
	scripts=['scripts/cstorage'],
	long_description=desc,
	test_suite="test"
)
