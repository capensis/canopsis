'use strict'
const commandLineCommands = require('../')

const cli = commandLineCommands([
  { name: 'help' },
  { name: 'run', definitions: [ { name: 'why', type: String } ] }
])

const command = cli.parse()

switch (command.name) {
  case 'help':
    console.log("I can't help you.")
    break
  case 'run':
    console.log(`${command.options.why}: this is not a good reason.`)
    break
  default:
    console.log('Unknown command.')
}
