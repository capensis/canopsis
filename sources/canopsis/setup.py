from setuptools import setup, find_packages

import canopsis

# TODO : if dependency is required
packages = ['canopsis', 'cconfiguration/cconfiguration',
	'cconfiguration/cconfiguration/manager',
	'ctimeserie/ctimeserie']

install_requires = [
]

with open('README') as f:
	desc = f.read()

setup(
	name=canopsis.__name__,
	version=canopsis.__version__,
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("Hypervisor"),
	license="AGPL v3",
	zip_safe=False,
	keywords="canopsis hypervision hypervisor monitoring",
	install_requires=install_requires,
	url="http://www.canopsis.org",
	packages=packages,
	long_description=desc,
	scripts=['scripts/canopsis']
)
