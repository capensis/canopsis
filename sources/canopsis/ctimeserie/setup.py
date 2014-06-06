from setuptools import setup

import ctimeserie

with open('README') as f:
	desc = f.read()

setup(
	name=ctimeserie.__name__,
	version=ctimeserie.__version__,
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("Performance data"),
	license="AGPL v3",
	zip_safe=False,
	keywords="perfdata performance canopsis",
	url="http://www.canopsis.org",
	packages=[ctimeserie.__name__],
	scripts=['scripts/ctimeserie'],
	long_description=desc,
	test_suite="test"
)
