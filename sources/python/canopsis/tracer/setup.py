from canopsis.common.setup import setup

install_requires = [
    'canopsis.common',
    'canopsis.confng',
    'canopsis.middleware'
]

setup(
    description='Canopsis tracer library',
    install_requires=install_requires,
    keywords='tracer')