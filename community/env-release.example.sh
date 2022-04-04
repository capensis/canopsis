export CANOPSIS_TAG="${tag}"

# all pour tout construire
#export CANOPSIS_DISTRIBUTION="all"
export CANOPSIS_DISTRIBUTION="debian-9"
export CANOPSIS_ENV_CONFIRM=1

# chemin vers la racine du dépôt canopsis core à utiliser. Cette varible n’est utilisée que par les scripts pour CAT.
export CANOPSIS_CORE_PATH="${HOME}/path/to/canopsis"

# Dans le cas où la version CANOPSIS_TAG n’est pas compatible avec le système de construction de packages, comme les DEB :
#export CANOPSIS_PACKAGE_TAG=<numéro de version>
# Dans le cas où la version DU PAQUET doit être changée (version 2, 3… rc1 etc) :
#export CANOPSIS_PACKAGE_REL=<numéro de version du paquet>

