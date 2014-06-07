from setuptools import setup, find_packages

with open('README') as f:
    desc = f.read()

install_requires = [
    'cconfiguration'
]

setup(
    name='ctimeserie',
    version='0.1',
    author="Capensis",
    author_email="canopsis@capensis.fr",
    description=("Performance data"),
    license="AGPL v3",
    zip_safe=False,
    keywords="perfdata performance canopsis",
    install_requires=install_requires,
    url="http://www.canopsis.org",
    packages=find_packages(exclude='test'),
    scripts=['scripts/ctimeserie'],
    long_description=desc,
    test_suite="test"
)
