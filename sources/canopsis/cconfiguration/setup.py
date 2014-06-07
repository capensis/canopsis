from setuptools import setup, find_packages

install_requires = [
    'canopsis'
]

with open('README') as f:
    desc = f.read()

setup(
    name="cconfiguration",
    version="0.1",
    author="Capensis",
    author_email="canopsis@capensis.fr",
    description=("Canopsis configuration"),
    license="AGPL v3",
    zip_safe=False,
    keywords="configuration canopsis",
    install_requires=install_requires,
    url = "http://www.canopsis.org",
    packages=find_packages(exclude=['test']),
    scripts=['scripts/cconfiguration'],
    long_description=desc,
    test_suite="test"
)
