#!/usr/bin/env node
'use strict'
const commandLineCommands = require('../')

const commands = [
  {
    name: 'eat',
    definitions: [
      { name: 'food', type: String, description: 'name of food' }
    ]
  },
  {
    name: 'sleep',
    definitions: [
      { name: 'hours', type: Number, description: 'Number of hours sleep' }
    ]
  }
]

const cli = commandLineCommands(commands)

console.log(cli.parse())
