from setuptools import setup

install_requires = [
	"pymongo", "ctimeserie", "cstorage"
]

setup(
	name="cmongo",
	version="0.1",
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("Mongo for canopsis"),
	license="AGPL v3",
	zip_safe=False,
	keywords="store performance mongo canopsis",
	install_requires=install_requires,
	url = "http://www.canopsis.org",
	packages=['cmongo'],
	scripts=['scripts/cmongo'],
	long_description=open('README').read(),
	test_suite="test"
)
