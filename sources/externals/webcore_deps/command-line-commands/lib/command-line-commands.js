'use strict'
const commandLineArgs = require('command-line-args')
const arrayify = require('array-back')

/**
 * Add a git-like command interface to your app. Wraps [command-line-args](https://github.com/75lb/command-line-args).
 *
 * @module command-line-commands
 * @example
 * ```js
 * const commandLineCommands = require('command-line-commands')
 *
 * // define your commands
 * const cli = commandLineCommands([
 *   { name: 'help' },
 *   { name: 'run', definitions: [ { name: 'why', type: String } ] }
 * ])
 *
 * // parse the command line
 * const command = cli.parse()
 *
 * // respond
 * switch (command.name) {
 *   case 'help':
 *     console.log("I can't help you.")
 *     break
 *   case 'run':
 *     console.log(`${command.options.why}: this is not a good reason.`)
 *     break
 *   default:
 *     console.log('Unknown command.')
 * }
 * ```
 *
 * Output (assumes your app name is `example`):
 * ```
 * $ example help
 * I can't help you.
 *
 * $ example run --why terror
 * terror: this is not a good reason.
 *
 * $ example hide
 * Unknown command.
 * ```
 */
module.exports = factory

/**
 * @class
 * @alias module:command-line-commands
 * @param commands {array}
 */
class CommandLineCommands {
  constructor (commands) {
    this.commands = commands
  }

  /**
   * @param [argv] {array}
   * @returns {object}
   */
  parse (argv) {
    if (argv) {
      argv = arrayify(argv)
    } else {
      /* if no argv supplied, assume we are parsing process.argv */
      argv = process.argv.slice(0)
      argv.splice(0, 2)
    }
    const commandName = argv.shift()
    const output = {}

    const commandDefinition = this.commands.find(c => c.name === commandName)
    if (commandDefinition) {
      const cli = commandLineArgs(commandDefinition.definitions)
      output.name = commandName
      output.options = cli.parse(argv)
    }
    return output
  }
}

function factory (commands) {
  return new CommandLineCommands(commands)
}
