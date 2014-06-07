from setuptools import setup, find_packages

install_requires = [
    "pymongo", "cstorage"
]

with open('README') as f:
    desc = f.read()

setup(
    name='cmongo',
    version='0.1',
    author="Capensis",
    author_email="canopsis@capensis.fr",
    description=("Mongo for canopsis"),
    license="AGPL v3",
    zip_safe=False,
    keywords="store mongo canopsis",
    install_requires=install_requires,
    url = "http://www.canopsis.org",
    packages=find_packages(exclude='test'),
    scripts=['scripts/cmongo'],
    long_description=desc,
    test_suite="test"
)
