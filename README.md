# Querybuilder Canopsis Brick

## Description

Query builder editor

## Screenshots



## Installation

You need to clone the git repository and copy directory to Canopsis path

    $ git clone https://git.canopsis.net/canopsis-ui-bricks/brick-querybuilder.git
    $ cp -r brick-querybuilder $CANOPSIS_PATH/var/www/canopsis

Then, you need to import specific schemas

    $ su - canopsis
    $ cp $CANOPSIS_PATH/var/www/canopsis/brick-querybuilder/schemas/* $CANOPSIS_PATH/etc/schema.d
    $ schema2db update

Then, you need to enable the brick

    $ su - canopsis
    $ webmodulemanager enable brick-querybuilder

You can see enabled bricks

    $ su - canopsis
    $ webmodulemanager list
    [u'core', u'uibase', u'monitoring', ..., **u'brick-querybuilder'**]

## Usage

See [Howto](https://git.canopsis.net/canopsis-ui-bricks/brick-querybuilder/blob/master/doc/index.rst)
