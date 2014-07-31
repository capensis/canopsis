from setuptools import setup

install_requires = [
	"pymongo",
]

setup(
	name = "pyperfstore2",
	version = "0.1",
	author = "Capensis",
	author_email = "canopsis@capensis.fr",
	description = ("Store performance data"),
	license = "AGPL v3",
	zip_safe = False,
	keywords = "nagios perfdata storage store performance canopsis",
	install_requires=install_requires,
	url = "http://www.canopsis.org",
	packages=['pyperfstore2'],
	scripts=['scripts/pyperfstore2'],
	data_files=[('~/etc/tasks.d', ['etc/tasks.d/task_pyperfstore.py'])],
)
