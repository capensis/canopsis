from setuptools import setup

install_requires = [
	"pymongo",
]

setup(
	name="cstorage",
	version="0.1",
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("Store data"),
	license="AGPL v3",
	zip_safe=False,
	keywords="storage store canopsis",
	install_requires=install_requires,
	url="http://www.canopsis.org",
	packages=['cstorage'],
	scripts=['scripts/cstorage'],
	long_description=open('README').read(),
	test_suite="test"
)
