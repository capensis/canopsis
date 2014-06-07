from setuptools import setup, find_packages

install_requires = [
]

with open('README') as f:
    desc = f.read()

setup(
    name='canopsis',
    version='0.1',
    author="Capensis",
    author_email="canopsis@capensis.fr",
    description=("Hypervisor"),
    license="AGPL v3",
    zip_safe=False,
    keywords="canopsis hypervision hypervisor monitoring",
    install_requires=install_requires,
    url="http://www.canopsis.org",
    packages=find_packages(exclude=['test']),
    long_description=desc,
    scripts=['scripts/canopsis']
)
