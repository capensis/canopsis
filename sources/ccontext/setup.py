from setuptools import setup

install_requires = [
]

setup(
	name="ccontext",
	version="0.1",
	author="Capensis",
	author_email="canopsis@capensis.fr",
	description=("Store ccontext"),
	license="AGPL v3",
	zip_safe=False,
	keywords="ccontext storage store canopsis ccontext",
	install_requires=install_requires,
	url="http://www.canopsis.org",
	packages=['ccontext'],
	scripts=['scripts/ccontext'],
	long_description=open('README').read(),
	test_suite="test"
)
