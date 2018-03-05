'use strict';

var _createClass = (function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; })();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var commandLineArgs = require('command-line-args');
var arrayify = require('array-back');

module.exports = factory;

var CommandLineCommands = (function () {
  function CommandLineCommands(commands) {
    _classCallCheck(this, CommandLineCommands);

    this.commands = commands;
  }

  _createClass(CommandLineCommands, [{
    key: 'parse',
    value: function parse(argv) {
      if (argv) {
        argv = arrayify(argv);
      } else {
        argv = process.argv.slice(0);
        argv.splice(0, 2);
      }
      var commandName = argv.shift();
      var output = {};

      var commandDefinition = this.commands.find(function (c) {
        return c.name === commandName;
      });
      if (commandDefinition) {
        var cli = commandLineArgs(commandDefinition.definitions);
        output.name = commandName;
        output.options = cli.parse(argv);
      }
      return output;
    }
  }]);

  return CommandLineCommands;
})();

function factory(commands) {
  return new CommandLineCommands(commands);
}