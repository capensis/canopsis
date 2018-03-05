'use strict';

var chalk = require('chalk'),
    optionPattern = /^(--([\w-]+)|-(\w+))(=(.*))?/,
    helpPattern = /(\s|^)(-h|--help)(\s|$)/,
    valuePattern = /^[^-].*/,
    useColors = true;

/* =============================================================================================
   PUBLIC APIs
   ============================================================================================= */

/**
 * Constructor for an option.
 * @param {String}  name             Option's name
 * @param {Object}  props            Option's properties
 * @param {String}  [prop.abbr]      Option's abbreviation
 * @param {String}  [prop.desc]      Option's description
 * @param {String}  [prop.valueName] Name of the option's value displayed in the usage message
 * @param {Boolean} [prop.flag]      Define if option is a flag
 * @param {*}       [prop.default]   Default value for the option
 * @constructor
 */
function Option(name, props) {
    this.name = name;
    this.abbr = props.abbr;
    this.valueName = props.valueName;
    this.desc = props.desc || '';
    this.flag = props.flag || false;
    this.defaultValue = props.default;
}

/**
 * Return pre-formatted usage information.
 * @return {Object}
 */
Option.prototype.getUsage = function() {
    var string = [];

    /* istanbul ignore else */
    if (this.abbr) {
        string.push('-' + this.abbr);

        if (!this.flag && this.valueName) {
            string.push(this.valueName);
        }

        string[string.length - 1] += ',';
    }

    string.push('--' + this.name);

    if (!this.flag && this.valueName) {
        string.push(this.valueName);
    }

    return {
        name: string.join(' '),
        desc: this.desc + (this.defaultValue ? (this.desc ? ' ' : '') + '[' + this.defaultValue + ']' : '')
    };
};



/**
 * Constructor for a command.
 * @param {String} name             Command's name
 * @param {Object} props            Command's properties
 * @param {String} [props.desc]     Command's description
 * @param {String} [props.callback] Callback function for the command
 * @constructor
 */
function Command(name, props) {
    this.script = null;
    this.options = {};
    this.commands = {};
    this.parsed = {};
    this.args = [];
    this.parent = null;

    this.name = name;
    this.desc = props && props.desc || '';
    this.callback = props && props.callback;
}

/**
 * Add a new option to the current command or the CLI's root. Following option formats are recognized:
 * - `-f [value]`
 * - `--foo [value]`
 * - `-f=value`
 * - `--foo=value`
 * @param  {String}  name             Option's name
 * @param  {Object}  props            Option's properties
 * @param  {String}  [prop.abbr]      Option's abbreviation
 * @param  {String}  [prop.desc]      Option's description
 * @param  {String}  [prop.valueName] Name of the option's value displayed in the usage message
 * @param  {Boolean} [prop.flag]      Define if option is a flag
 * @param  {*}       [prop.default]   Default value for the option
 * @return {Command}
 */
Command.prototype.option = function(name, props) {
    var option = new Option(name, props);

    this.options[name] = option;

    return this;
};

/**
 * Add a new (sub-)command to the current command or CLI's root.
 * @param  {String} name             Command's name
 * @param  {Object} props            Command's properties
 * @param  {String} [props.desc]     Command's description
 * @param  {String} [props.callback] Command's callback function executed if the command is run
 * @return {Command} (Sub-)command instance
 */
Command.prototype.command = function(name, props) {
    var command = new Command(name, props);

    command.parent = this;

    this.commands[name] = command;

    return command;
};

/**
 * Set the name of the script's executable.
 * @param  {String} name Name of the executable
 * @return {Command}
 */
Command.prototype.scriptName = function(name) {
    this.script = name;

    return this;
};

/**
 * Disable coloring in usage and error messages
 * @return {Command}
 */
Command.prototype.noColors = function() {
    useColors = false;

    return this;
};

/**
 * Parse the command line arguments
 * @param  {Array.<String>} [argv] Array of arguments
 * @return {Object} List of parsed arguments
 */
Command.prototype.parse = function(argv) {
    var expectCommand = Object.keys(this.commands).length > 0,
        command,
        parsed;

    /* istanbul ignore next */
    argv = argv || process.argv.slice(2);

    // contains sub-commands
    if (expectCommand) {
        // no argument or just an option
        if (!argv[0] || argv[0].charAt(0) === '-') {
            // show usage information
            if (helpPattern.test(argv.join(' '))) {
                return this.usage();
            }

            // call the callback with parsed arguments
            if (typeof this.callback == 'function') {
                this._parseArgv(argv);
                parsed = this._getParsed();
                this.callback(parsed);

                return parsed;
            }

            return this._printError('Missing command for "' + this._getScriptName() + '".');
        }

        // handle a command
        if ((command = this.commands[argv[0]])) {
            return command.parse.call(command, argv.slice(1));
        }

        return this._printError('Unknown command "' + argv[0] + '".');
    }

    this._parseArgv(argv);
    parsed = this._getParsed();

    if (typeof this.callback == 'function') {
        this.callback(parsed);
    }

    return parsed;
};

/**
 * End modifying current command and return to the parent
 * @return {Command} Parent command (if any)
 */
Command.prototype.end = function() {
    return this.parent || this;
};

/**
 * Print command's usage message on the STDOUT.
 */
Command.prototype.usage = function() {
    var commandNames,
        optionNames,
        longest,
        options;

    function spaces(num) {
        return new Array(num + 3).join(' ');
    }

    commandNames = Object.keys(this.commands);
    options = this._getOptions();
    optionNames = Object.keys(options);

    // print usage line
    write(chalk.bold('\nUsage:'));
    write(' ' + this._getCommandChain().join(' '));

    if (commandNames.length) {
        write(yellow(' <command>'));
    }

    if (optionNames.length) {
        write(cyan(' [options]'));
    }

    write('\n');

    // print commands list
    if (commandNames.length) {
        longest = commandNames.reduce(function(a, b) {
            return a.length > b.length ? a : b;
        }).length;

        writeLine(chalk.bold(yellow('Commands:\n')));
        commandNames.sort().forEach(function(name) {
            var command = this.commands[name];

            write('  ' + command.name);
            write(spaces(longest - command.name.length));
            write(grey(command.desc) + '\n');
        }, this);
    }

    // print options list
    if (optionNames.length) {
        optionNames.forEach(function(name) {
            options[name] = options[name].getUsage();
        });

        longest = optionNames.reduce(function(a, b) {
            return options[a].name.length > options[b].name.length ? a : b;
        });
        longest = options[longest].name.length;

        writeLine(chalk.bold(cyan('Options:\n')));
        optionNames.sort().forEach(function(name) {
            var option = options[name];

            write('  ' + option.name);
            write(spaces(longest - option.name.length));
            write(grey(option.desc) + '\n');
        });
    }

    write('\n');
};

/**
 * Resets all properties of the command
 * @return {Command}
 */
Command.prototype.reset = function() {
    this.script = null;
    this.options = {};
    this.commands = {};
    this.parsed = {};
    this.args = [];
    this.parent = null;

    return this;
};

/* =============================================================================================
   PRIVATE APIs
   ============================================================================================= */

/**
 * Get an option an it's value based on a given argument.
 * @param  {String} arg Command line argument
 * @return {Object}
 * @private
 */
Command.prototype._getOption = function(arg) {
    var options = this._getOptions(),
        option,
        match,
        value,
        name,
        o;

    // matches --foo, -f and --foo=<value>
    if ((match = optionPattern.exec(arg))) {
        name = match[2] || match[3];
        value = match[5];

        // handle -h / --help and stop execution
        if (name === 'help' || name === 'h') {
            this.usage();
            process.exit(0);
        }

        for (o in options) {
            if ((option = options[o]) !== undefined &&
                name === o || name === option.abbr) {
                return {
                    name: option.name,
                    defaultValue: option.defaultValue,
                    value: option.flag || value
                };
            }
        }

        return {
            name: name,
            value: value,
            unknown: true
        };
    }

    return null;
};

/**
 * Write an error message to the STDOUT.
 * @param {String} message Error message to print
 * @private
 */
Command.prototype._printError = function(message) {
    writeLine(red(chalk.bold('Error: ') + message));
    this.usage();
    process.exit(1);
};

/**
 * Get the name of the script/command.
 * @return {String}
 * @private
 */
Command.prototype._getScriptName = function() {
    return this.name || this.script;
};

/**
 * Get all available options, merge parent's options if any.
 * @return {Object}
 * @private
 */
Command.prototype._getOptions = function() {
    var result = {},
        parents,
        name;

    for (name in this.options) {
        /* istanbul ignore else */
        if (this.options[name] !== undefined) {
            result[name] = this.options[name];
        }
    }

    if (this.parent) {
        parents = this.parent._getOptions();
        for (name in parents) {
            /* istanbul ignore else */
            if (parents[name] !== undefined) {
                result[name] = parents[name];
            }
        }
    }

    return result;
};

/**
 * Get all parsed options and arguments.
 * @return {Object}
 * @private
 */
Command.prototype._getParsed = function() {
    var result = {},
        parents,
        item,
        name;

    for (name in this.parsed) {
        /* istanbul ignore else */
        if ((item = this.parsed[name]) !== undefined) {
            result[name] = item;
        }
    }

    for (name in this.options) {
        /* istanbul ignore else */
        if ((item = this.options[name]) !== undefined &&
            !result[name] && item.defaultValue) {
            result[name] = item.defaultValue;
        }
    }

    if (this.parent) {
        parents = this.parent._getParsed();
        for (name in parents) {
            /* istanbul ignore else */
            if ((item = parents[name]) !== undefined && !result[name]) {
                result[name] = item;
            }
        }
    }

    /* istanbul ignore else */
    if (this.args) {
        this.args.forEach(function(arg, index) {
            result[index] = arg;
        });
    }

    return result;
};

/**
 * Parse given arguments as options/values
 * @param {Array.<String>} argv Array of arguments
 * @private
 */
Command.prototype._parseArgv = function(argv) {
    var option,
        arg,
        i;

    for (i = 0; i < argv.length; i++) {
        arg = argv[i];

        // it's an option
        if ((option = this._getOption(arg))) {
            if (!option.value) {
                // take the value from the next argument
                if (argv[i + 1] && valuePattern.test(argv[i + 1])) {
                    option.value = argv[i + 1];
                    i++;
                } else if (option.unknown) {
                    option.value = true;
                } else {
                    this._printError('Missing value for "' + option.name + '" option.');
                }
            }

            this.parsed[option.name] = option.value;
        } else {
            this.args.push(arg);
        }
    }
};

/**
 * Return the script - command - sub-command(s) chain of the current command
 * @return {Array.<String>}
 * @private
 */
Command.prototype._getCommandChain = function() {
    var result = [this.script || this.name];

    return this.parent ? this.parent._getCommandChain().concat(result) : result;
};

/**
 * Write a text to the STDOUT without adding new line character.
 * @param {String} text Text to write
 * @private
 */
/* istanbul ignore next */
function write(text) {
    process.stdout.write(text);
}

/**
 * Write a text to the STDOUT in a new line.
 * @param {String} text Text to write
 * @private
 */
/* istanbul ignore next */
function writeLine(text) {
    process.stdout.write('\n' + text + '\n');
}

function red(text) {
    return useColors ? chalk.red(text) : text;
}

function yellow(text) {
    return useColors ? chalk.yellow(text) : text;
}

function cyan(text) {
    return useColors ? chalk.cyan(text) : text;
}

function grey(text) {
    return useColors ? chalk.grey(text) : text;
}

module.exports = new Command();

module.exports.script = require('path').basename(process.argv[1]);
module.exports.Option = Option;
module.exports.Command = Command;
