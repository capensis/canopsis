from setuptools import setup, find_packages
import os

install_requires = [
    "pymongo", 'cconfiguration'
]

with open(os.path.join(os.path.dirname(__file__), 'README')) as f:
    desc = f.read()

setup(
    name='cstorage',
    version='0.1',
    author="Capensis",
    author_email="canopsis@capensis.fr",
    description=("Store data"),
    license="AGPL v3",
    zip_safe=False,
    keywords="storage store canopsis",
    install_requires=install_requires,
    url="http://www.canopsis.org",
    packages=find_packages(exclude='test'),
    long_description=desc,
    test_suite="test"
)
