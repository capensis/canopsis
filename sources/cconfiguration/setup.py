from setuptools import setup

install_requires = [
]

setup(
	name="cconfiguration",
	version="0.1",
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("Canopsis configuration"),
	license="AGPL v3",
	zip_safe=False,
	keywords="configuration canopsis",
	install_requires=install_requires,
	url = "http://www.canopsis.org",
	packages=['cconfiguration'],
	scripts=['scripts/cconfiguration'],
	long_description=open('README').read(),
	test_suite="test"
)
