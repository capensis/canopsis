from setuptools import setup

import ccontext

install_requires = [
]

with open('README').read():
	desc = f.read()

setup(
	name=ccontext.__name__,
	version=ccontext.__version__,
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("Store ccontext"),
	license="AGPL v3",
	zip_safe=False,
	keywords="ccontext storage store canopsis ccontext",
	install_requires=install_requires,
	url="http://www.canopsis.org",
	packages=[ccontext.__name__],
	scripts=['scripts/ccontext'],
	long_description=desc,
	test_suite="test"
)
