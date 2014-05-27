from setuptools import setup

install_requires = [
	"pymongo", "ctimeserie", "cstorage"
]

setup(
	name="cperfdata",
	version="0.1",
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("performance data"),
	license="AGPL v3",
	zip_safe=False,
	keywords="nagios perfdata storage store performance canopsis",
	install_requires=install_requires,
	url="http://www.canopsis.org",
	packages=['cperfdata'],
	scripts=['scripts/cperfdata'],
	long_description=open('README').read(),
	test_suite="test"
)
