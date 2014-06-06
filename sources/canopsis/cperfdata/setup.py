from setuptools import setup

import cperfdata

install_requires = [
]

with open('README') as f:
	desc = f.read()

setup(
	name=cperfdata.__name__,
	version=cperfdata.__version__,
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("performance data"),
	license="AGPL v3",
	zip_safe=False,
	keywords="nagios perfdata storage store performance canopsis",
	install_requires=install_requires,
	url="http://www.canopsis.org",
	packages=[cperfdata.__name__],
	scripts=['scripts/cperfdata'],
	long_description=desc,
	test_suite="test"
)
