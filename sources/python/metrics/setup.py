from canopsis.common.setup import setup

install_requires = [
    'canopsis.common',
]

setup(
    description='Canopsis metrics',
    install_requires=install_requires,
    test_suite='test',
    keywords='metrics'
)
