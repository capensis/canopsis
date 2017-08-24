from canopsis.common.setup import setup

install_requires = [
    'canopsis.common',
]

setup(
    description='Canopsis logger',
    install_requires=install_requires,
    test_suite='test',
    keywords='logger')
