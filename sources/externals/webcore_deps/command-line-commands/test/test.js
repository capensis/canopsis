'use strict'
var test = require('tape')
var commandLineCommands = require('../')

test('parse: simple', function (t) {
  var commands = [
    {
      name: 'eat',
      definitions: [ { name: 'food' } ]
    },
    {
      name: 'sleep',
      definitions: [ { name: 'hours' } ]
    }
  ]
  var cli = commandLineCommands(commands)
  var command = cli.parse([ 'eat', '--food', 'peas' ])
  t.deepEqual(command, {
    name: 'eat',
    options: { food: 'peas' }
  })
  command = cli.parse([ 'sleep', '--hours', '2' ])
  t.deepEqual(command, {
    name: 'sleep',
    options: { hours: '2' }
  })
  t.end()
})

test('parse: no definitions', function (t) {
  var commands = [ { name: 'eat' } ]
  var cli = commandLineCommands(commands)
  var command = cli.parse([ 'eat' ])
  t.deepEqual(command, {
    name: 'eat',
    options: { }
  })
  t.end()
})

test('parse: no definitions, but options passed', function (t) {
  var commands = [ { name: 'eat' } ]
  var cli = commandLineCommands(commands)
  t.throws(function () {
    cli.parse([ 'eat', '--food', 'peas' ])
  })
  t.end()
})
