from setuptools import setup

install_requires = [
	"cstorage",
]

setup(
	name="ctopology",
	version="0.1",
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("Store topology"),
	license="AGPL v3",
	zip_safe=False,
	keywords="ctopology storage store canopsis ctopology",
	install_requires=install_requires,
	url="http://www.canopsis.org",
	packages=['ctopology'],
	scripts=['scripts/ctopology'],
	long_description=open('README').read(),
	test_suite="test"
)
