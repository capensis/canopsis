from setuptools import setup

import cmongo

install_requires = [
	#"pymongo", "ctimeserie", "cstorage"
]

with open('README') as f:
	desc = f.read()

setup(
	name=cmongo.__name__,
	version=cmongo.__version__,
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("Mongo for canopsis"),
	license="AGPL v3",
	zip_safe=False,
	keywords="store mongo canopsis",
	install_requires=install_requires,
	url = "http://www.canopsis.org",
	packages=[cmongo.__name__],
	scripts=['scripts/cmongo'],
	long_description=desc,
	test_suite="test"
)
