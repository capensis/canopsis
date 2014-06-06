from setuptools import setup

import cconfiguration

install_requires = [
]

with open('README') as f:
	desc = f.read()

setup(
	name=cconfiguration.__name__,
	version=cconfiguration.__version__,
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("Canopsis configuration"),
	license="AGPL v3",
	keywords="configuration canopsis",
	install_requires=install_requires,
	url = "http://www.canopsis.org",
	packages=[cconfiguration.__name__],
	scripts=['scripts/cconfiguration'],
	long_description=desc,
	test_suite="test"
)
