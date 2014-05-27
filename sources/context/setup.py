from setuptools import setup

install_requires = [
	"pymongo",
]

setup(
	name = "context",
	version = "0.1",
	author = "Capensis",
	author_email = "canopsis@capensis.fr",
	description = ("Store context"),
	license = "AGPL v3",
	zip_safe = False,
	keywords = "context storage store canopsis context",
	install_requires=install_requires,
	url = "http://www.canopsis.org",
	packages=['context'],
	scripts=['scripts/context'],
	long_description=open('README').read(),
	test_suite="test"
)
