from setuptools import setup

setup(
	name="ctimeserie",
	version="0.1",
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("Performance data"),
	license="AGPL v3",
	zip_safe=False,
	keywords="perfdata performance canopsis",
	url="http://www.canopsis.org",
	packages=['ctimeserie'],
	scripts=['scripts/ctimeserie'],
	long_description=open('README').read(),
	test_suite="test"
)
