/* global describe, it, beforeEach, afterEach */
/*jshint -W030 */

'use strict';

var rewire = require('rewire'),
    sc = rewire('../'),
    expect = require('chai').expect;

describe('subcommander', function() {
    var oldExit,
        oldWrite,
        oldWriteLine,
        output,
        code;

    beforeEach(function() {
        oldExit = process.exit;
        oldWrite = sc.__get__('write');
        oldWriteLine = sc.__get__('writeLine');

        output = '';
        code = 0;

        process.exit = function(c) {
            code = c;
        };

        sc.__set__('write', function(text) {
            output += text;
        });

        sc.__set__('writeLine', function(text) {
            output += '\n' + text + '\n';
        });
    });

    afterEach(function() {
        process.exit = oldExit;
        sc.__set__('write', oldWrite);
        sc.__set__('writeLine', oldWriteLine);
        sc.reset();
    });

    it('should expose its public API', function() {
        expect(sc).to.be.instanceof(sc.Command);
        expect(Object.getPrototypeOf(sc)).to.contain.keys(
            ['option', 'command', 'parse', 'usage', 'scriptName', 'noColors', 'end']
        );
    });

    it('should parse undefined -x option as a flag', function() {
        expect(sc.parse(['-f'])).to.deep.equal({
            f: true
        });
    });

    it('should parse undefined --xxx option as a flag', function() {
        expect(sc.parse(['--foo'])).to.deep.equal({
            foo: true
        });
    });

    it('should parse undefined -x option and its value', function() {
        expect(sc.parse(['-f', 'bar'])).to.deep.equal({
            f: 'bar'
        });
    });

    it('should parse undefined --xxx option and its value', function() {
        expect(sc.parse(['--foo', 'bar'])).to.deep.equal({
            foo: 'bar'
        });
    });

    it('should parse undefined -x=yyy option and its value', function() {
        expect(sc.parse(['-f=bar'])).to.deep.equal({
            f: 'bar'
        });
    });

    it('should parse undefined --xxx=yyy option and its value', function() {
        expect(sc.parse(['--foo=bar'])).to.deep.equal({
            foo: 'bar'
        });
    });

    it('should add unrecognized arguments to the output', function() {
        expect(sc.parse(['--foo=foo', 'bar', '-f', 'quux', 'baz'])).to.deep.equal({
            foo: 'foo',
            f: 'quux',
            '0': 'bar',
            '1': 'baz'
        });
    });

    it('should define an option', function() {
        var expected = {
            foo: 'bar'
        };

        function redefine() {
            return sc.reset().option('foo', {
                abbr: 'f',
                desc: 'desc for foo'
            });
        }

        redefine();

        expect(sc.options).to.have.key('foo');
        expect(sc.options.foo).to.be.instanceof(sc.Option);
        expect(sc.parse(['-f', 'bar'])).to.deep.equal(expected);
        expect(redefine().parse(['--foo', 'bar'])).to.deep.equal(expected);
        expect(redefine().parse(['-f=bar'])).to.deep.equal(expected);
        expect(redefine().parse(['--foo=bar'])).to.deep.equal(expected);
    });

    describe('option', function() {
        it('should handle a flag', function() {
            var expected = {
                foo: true,
                '0': 'bar'
            };

            expect(
                sc.option('foo', {
                    abbr: 'f',
                    flag: true,
                    desc: 'desc for foo flag'
                }).parse(['-f', 'bar'])
            ).to.deep.equal(expected);

            expect(
                sc.reset().option('foo', {
                    abbr: 'f',
                    flag: true,
                    desc: 'desc for foo flag'
                }).parse(['--foo', 'bar'])
            ).to.deep.equal(expected);
        });

        it('should return a default value', function() {
            var expected = {
                foo: 'baz',
                '0': 'bar'
            };

            sc.option('foo', {
                abbr: 'f',
                desc: 'desc for foo',
                default: 'baz'
            });

            expect(sc.parse(['bar'])).to.deep.equal(expected);
        });

        it('should override a default value', function() {
            var expected = {
                foo: 'quux',
                '0': 'bar'
            };

            sc.option('foo', {
                abbr: 'f',
                desc: 'desc for foo',
                default: 'baz'
            });

            expect(sc.parse(['bar', '--foo', 'quux'])).to.deep.equal(expected);
        });

        it('should return a pre-formatted usage information', function() {
            var expected = {
                name: '-f value, --foo value',
                desc: 'desc for foo [bar]'
            };

            sc.option('foo', {
                abbr: 'f',
                desc: 'desc for foo',
                default: 'bar',
                valueName: 'value'
            });

            expect(sc.options.foo.getUsage()).to.deep.equal(expected);
        });

        it('should use empty string if no description given', function() {
            var expected = {
                name: '-f value, --foo value',
                desc: '[bar]'
            };

            sc.option('foo', {
                abbr: 'f',
                desc: '',
                default: 'bar',
                valueName: 'value'
            });

            expect(sc.options.foo.getUsage()).to.deep.equal(expected);
        });

        it('should throw an error if no value for option specified', function() {
            sc.option('foo', {
                abbr: 'f',
                desc: 'desc for foo'
            });

            sc.parse(['-f']);

            expect(code).to.equal(1);
            expect(output).to.equal(
                '\n\u001b[1mError: \u001b[22mMissing value for \"foo\" option.\n\u001b[1m\n' +
                'Usage:\u001b[22m  [options]\n\n\u001b[1mOptions:\n\u001b[22m\n  -f, --foo  desc for foo\n\n'
            );
        });
    });

    it('should define a command', function() {
        sc.command('foo', {
            desc: 'desc for foo',
            callback: function() {}
        });

        expect(sc.commands).to.have.key('foo');
        expect(sc.commands.foo).to.be.instanceof(sc.Command);
    });

    it('should report an error if no command was given', function() {
        sc.command('foo', {
            desc: 'desc for foo',
            callback: function() {}
        });

        sc.parse(['--bar', 'baz']);

        expect(code).to.equal(1);
        expect(output).to.equal(
            '\n\u001b[31m\u001b[1mError: \u001b[22mMissing command for \"null\".\u001b[39m\n' +
            '\u001b[1m\nUsage:\u001b[22m \u001b[33m <command>\u001b[39m\n\n' +
            '\u001b[1m\u001b[33mCommands:\n\u001b[39m\u001b[22m\n  foo  \u001b[90mdesc for foo\u001b[39m\n\n'
        );
    });

    describe('command', function() {
        it('should execute its callback with parsed arguments', function(done) {
            sc.command('foo', {
                desc: 'desc for foo',
                callback: function(parsed) {
                    expect(parsed).to.deep.equal({
                        bar: 'baz',
                        '0': 'quux'
                    });

                    done();
                }
            });

            sc.parse(['foo', '--bar', 'baz', 'quux']);
        });

        it('should define a sub-command', function() {
            var foo = sc.command('foo', {
                    desc: 'desc for foo',
                    callback: function() {}
                }),
                bar = foo.command('bar', {
                    desc: 'desc for bar',
                    callback: function() {}
                });

            expect(sc.commands).to.have.key('foo');
            expect(sc.commands.foo).to.equal(foo);

            expect(sc.commands.foo.commands).to.have.key('bar');
            expect(sc.commands.foo.commands.bar).to.equal(bar);
        });

        it('should execute its callback with parsed arguments', function(done) {
            var returned;

            sc.command('foo', {
                desc: 'desc for foo',
                callback: function(parsed) {
                    process.nextTick(function() {
                        expect(returned).to.deep.equal(parsed);

                        done();
                    });
                }
            });

            returned = sc.parse(['foo', '--bar', 'baz', 'quux']);
        });

        it('should allow commands containing dashes', function(done) {
            sc.command('has-dash', {
                desc: 'contains a dash in its name',
                callback: function(parsed) {
                    expect(parsed).to.deep.equal({
                        f: 'foo',
                        bar: 'bar'
                    });

                    done();
                }
            });

            sc.parse(['has-dash', '-f', 'foo', '--bar', 'bar']);
        });

        it('should execute its callback even if it has sub-commands defined', function(done) {
            var returned;

            sc.callback = function(parsed) {
                process.nextTick(function() {
                    expect(returned).to.deep.equal(parsed);
                    done();
                });
            };

            sc
                .option('version', {
                    abbr: 'v',
                    desc: 'desc for version',
                    flag: true
                })
                .command('foo', {
                    desc: 'desc for foo',
                    callback: function() {
                        throw new Error('You should never call me!');
                    }
                });

            returned = sc.parse(['--version', 'foo', 'bar']);
        });
    });

    describe('sub-command', function() {
        it('should execute its callback with parsed arguments', function(done) {
            sc
                .command('foo', {
                    desc: 'desc for foo'
                })
                .command('bar', {
                    desc: 'desc for bar',
                    callback: function(parsed) {
                        expect(parsed).to.deep.equal({
                            'baz': 'quux'
                        });

                        done();
                    }
                });

            sc.parse(['foo', 'bar', '--baz', 'quux']);
        });

        it('should return its parent using end() method', function() {
            sc
                .command('foo', {
                    desc: 'desc for foo',
                    callback: function() {}
                })
                .end()
                .command('bar', {
                    desc: 'desc for bar',
                    callback: function() {}
                })
                .end()
                .end();

            expect(sc.commands).to.have.keys(['foo', 'bar']);
            expect(sc.commands.foo.commands).to.be.empty;
        });

        it('should inherit its parent\'s options and its default values', function(done) {
            function complete(parsed) {
                expect(parsed).to.deep.equal({
                    'baz': 'quux',
                    'quux': 'quax'
                });

                done();
            }

            sc
                .command('foo', {
                    desc: 'desc for foo'
                })
                .option('baz', {
                    abbr: 'b',
                    desc: 'desc for baz',
                    default: 'quux'
                })
                .command('bar', {
                    desc: 'desc for bar',
                    callback: complete
                })
                .option('quux', {
                    abbr: 'q',
                    desc: 'desc for quux'
                });

            sc.parse(['foo', 'bar', '--quux', 'quax']);
        });

        it('should throw an error on unknown command', function() {
            sc
                .command('foo', {
                    desc: 'desc for foo',
                    callback: function() {}
                });

            sc.parse(['bar']);

            expect(code).to.equal(1);
            expect(output).to.equal(
                '\n\u001b[1mError: \u001b[22mUnknown command \"bar\".\n\u001b[1m' +
                '\nUsage:\u001b[22m  <command>\n\n\u001b[1mCommands:\n\u001b[22m\n  foo  desc for foo\n\n'
            );
        });
    });

    it('should build usage information', function() {
        sc
            .option('baz', {
                abbr: 'b',
                desc: 'desc for baz',
                valueName: 'value',
                default: 'quux'
            })
            .option('quux', {
                abbr: 'q',
                desc: 'desc for quux',
                flag: true
            })
            .command('foo', {
                desc: 'desc for foo',
                callback: function() {}
            })
            .end()
            .command('bar', {
                desc: 'desc for bar',
                callback: function() {}
            });

        sc.usage();

        expect(output).to.equal(
            '\u001b[1m\nUsage:\u001b[22m \u001b[33m <command>\u001b[39m\u001b[36m [options]\u001b[39m\n\n' +
            '\u001b[1m\u001b[33mCommands:\n\u001b[39m\u001b[22m\n' +
            '  bar  \u001b[90mdesc for bar\u001b[39m\n  foo  \u001b[90mdesc for foo\u001b[39m\n\n' +
            '\u001b[1m\u001b[36mOptions:\n\u001b[39m\u001b[22m\n' +
            '  -b value, --baz value  \u001b[90mdesc for baz [quux]\u001b[39m\n' +
            '  -q, --quux             \u001b[90mdesc for quux\u001b[39m\n\n'
        );
    });

    it('should build usage information for a sub-command', function() {
        sc
            .option('baz', {
                abbr: 'b',
                desc: 'desc for baz',
                valueName: 'value',
                default: 'quux'
            })
            .option('quux', {
                abbr: 'q',
                desc: 'desc for quux',
                flag: true
            })
            .command('foo', {
                desc: 'desc for foo',
                callback: function() {}
            })
            .command('bar', {
                desc: 'desc for bar',
                callback: function() {}
            });

        sc.parse(['foo', 'bar', '-h']);

        expect(output).to.equal(
            '\u001b[1m\nUsage:\u001b[22m  foo bar\u001b[36m [options]\u001b[39m\n\n' +
            '\u001b[1m\u001b[36mOptions:\n\u001b[39m\u001b[22m\n' +
            '  -b value, --baz value  \u001b[90mdesc for baz [quux]\u001b[39m\n' +
            '  -q, --quux             \u001b[90mdesc for quux\u001b[39m\n\n'
        );
    });

    it('should print usage information if -h / --help flag passed', function() {
        sc
            .option('baz', {
                abbr: 'b',
                desc: 'desc for baz',
                valueName: 'value',
                default: 'quux'
            })
            .option('quux', {
                abbr: 'q',
                desc: 'desc for quux',
                flag: true
            })
            .command('foo', {
                desc: 'desc for foo',
                callback: function() {}
            })
            .end()
            .command('bar', {
                desc: 'desc for bar',
                callback: function() {}
            });

        sc.parse(['-h']);

        expect(output).to.equal(
            '\u001b[1m\nUsage:\u001b[22m \u001b[33m <command>\u001b[39m\u001b[36m [options]\u001b[39m\n\n' +
            '\u001b[1m\u001b[33mCommands:\n\u001b[39m\u001b[22m\n  bar  \u001b[90mdesc for bar\u001b[39m\n' +
            '  foo  \u001b[90mdesc for foo\u001b[39m\n\n\u001b[1m\u001b[36mOptions:\n\u001b[39m\u001b[22m\n' +
            '  -b value, --baz value  \u001b[90mdesc for baz [quux]\u001b[39m\n' +
            '  -q, --quux             \u001b[90mdesc for quux\u001b[39m\n\n'
        );
    });

    it('should print usage information if -h / --help flag passed 2', function() {
        sc
            .option('baz', {
                abbr: 'b',
                desc: 'desc for baz',
                valueName: 'value',
                default: 'quux'
            })
            .option('quux', {
                abbr: 'q',
                desc: 'desc for quux',
                flag: true
            });

        sc.parse(['-h']);

        expect(output).to.equal(
            '\u001b[1m\nUsage:\u001b[22m \u001b[36m [options]\u001b[39m\n\n' +
            '\u001b[1m\u001b[36mOptions:\n\u001b[39m\u001b[22m\n' +
            '  -b value, --baz value  \u001b[90mdesc for baz [quux]\u001b[39m\n' +
            '  -q, --quux             \u001b[90mdesc for quux\u001b[39m\n\n'
        );
    });

    it('should set script name', function() {
        sc.scriptName('foo');

        sc.usage();

        expect(output).to.equal('\u001b[1m\nUsage:\u001b[22m foo\n\n');
    });

    it('should allow to disable colors', function() {
        sc
            .noColors()
            .option('baz', {
                abbr: 'b',
                desc: 'desc for baz',
                valueName: 'value',
                default: 'quux'
            })
            .command('foo', {
                desc: 'desc for foo',
                callback: function() {}
            })
            .end()
            .command('bar', {
                desc: 'desc for bar',
                callback: function() {}
            });

        sc.usage();

        expect(output).to.equal(
            '\u001b[1m\nUsage:\u001b[22m  <command> [options]\n\n\u001b[1mCommands:\n' +
            '\u001b[22m\n  bar  desc for bar\n  foo  desc for foo\n\n\u001b[1mOptions:\n' +
            '\u001b[22m\n  -b value, --baz value  desc for baz [quux]\n\n'
        );
    });
});
