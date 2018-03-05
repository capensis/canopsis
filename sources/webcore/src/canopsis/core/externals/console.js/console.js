var _baseConsole = console;

//define Array.contains if not defined before
if(Array.prototype.contains === undefined) {
    Array.prototype.contains = function(obj) {
        var i = this.length;
        while (i--) {
            if (this[i] === obj) {
                return true;
            }
        }
        return false;
    };
}

console = {
    _filter: undefined,

    // internal utilities 
    internal: {
        generateMessageAdditions: function(consoleObject, argumentsArray, forceDisplay) {
            var args = [];
            var i;

            //getting author name
            var file_split;
            try {
                file_split = new Error().stack.split('\n')[3].split('/');
            } catch(e) {
                //console not supported, exiting.
                //TODO log the message with the standard output
                // alert(argumentsArray);
                return null;
            }

            //consoleObject._baseConsole.log("file split", file_split);

            var file_location = file_split[file_split.length - 1];
            for(i = 1, l = consoleObject.tags._authorFoldersDisplay; i < l; i++) {
                file_location = file_split[file_split.length - i - 1] + "/" + file_location;
            }

            file_location = file_location.replace(')','');

            var file_location_split = file_location.split(':');

            var filename = file_location_split[0];
            var line_number = file_location_split[1] + ":" + file_location_split[2];

            var selectedTags = console.tags._selectedTags;
            var selectedTagsContainsFilename = (selectedTags ? selectedTags.contains(filename): false);

            if (selectedTags !== undefined && selectedTags !== null && selectedTagsContainsFilename) {
                return null;
            }

            var args_tags = "[" + filename + "][" + line_number + "]";

            var tagMatchSelectedTag = false;
            if(consoleObject.tags._tags !== undefined) {
                for (i = 0, l = consoleObject.tags._tags.length - 1; i <= l ; i++) {
                    if (selectedTags.contains(consoleObject.tags._tags[i])) {
                        tagMatchSelectedTag = true;
                    }
                    args_tags += "["+ consoleObject.tags._tags[i] +"]";
                }
            }

            if(! forceDisplay) {
                if(consoleObject.tags._muteAllByDefault === true) {
                    if (selectedTags !== undefined && selectedTags !== null && !selectedTagsContainsFilename) {
                            return null;
                    }
                    if(tagMatchSelectedTag === false) {
                            return null;
                    }
                } else {
                    if (selectedTags !== undefined && selectedTags !== null && selectedTagsContainsFilename) {
                            return null;
                    }
                    if(tagMatchSelectedTag === true) {
                            return null;
                    }
                }
            }

            if(!! consoleObject.style._colors) {
                args_tags = "%c" + args_tags;
                args.push(args_tags);
                args.push('background: #444; color: #eee; border-radius:4px;padding:2px');
            } else {
                args.push(args_tags);
            }

            if(consoleObject._filter !== "" && consoleObject._filter !== undefined && consoleObject._filter !== null) {
                for (i = 0, l = args.length; i < l; i++) {
                    if(typeof args[i] === "string") {
                        var regex = new RegExp(consoleObject._filter);
                        if(regex.test(args[i])) {
                            if(! forceDisplay)
                                return null;
                        }
                    }
                }
            }

            if(!! consoleObject.stacks.display) {
                var err = new Error();

                consoleObject.stacks._stacks.push(err.stack);
                args.push("[>" + consoleObject.stacks._stacks.length + "]");
            }

            for (i = 0, l = argumentsArray.length; i < l; i++) {
                args.push(argumentsArray[i]);
            }

            return args;
        }
    },

    // Ctor 

    init: function() {
        console.log("load settings from localStorage, if possible");
        this.settings.load();
        if(this.tags._selectedTags === undefined || this.tags._selectedTags === null) {
            this.tags._selectedTags = [];
        }
        this.tags.flush();
    },

    // Original console method wrappers 

    log: function() {
        var args = this.internal.generateMessageAdditions(this, arguments);

        if(args !== null) {
            _baseConsole.log.apply(_baseConsole, args);
            this.backends.send("log", args);
        }
    },

    group: function() {
        var args = console.internal.generateMessageAdditions(this, arguments, true);

        if(args !== null) {
            _baseConsole.group.apply(_baseConsole, args);
            this.backends.send("group", args);
        }
    },

    groupEnd: function() {
        var args = this.internal.generateMessageAdditions(this, arguments, true);

        if(args !== null) {
            _baseConsole.groupEnd.apply(_baseConsole, args);
            this.backends.send("groupEnd", args);
        }
    },

    groupCollapsed: function() {
        var args = this.internal.generateMessageAdditions(this, arguments, true);

        if(args !== null) {
            _baseConsole.groupCollapsed.apply(_baseConsole, args);
            this.backends.send("groupCollapsed", args);
        }
    },

    info: function() {
        var args = this.internal.generateMessageAdditions(this, arguments);

        if(args !== null) {
            _baseConsole.info.apply(_baseConsole, args);
            this.backends.send("info", args);
        }
    },

    warn: function() {
        var args = this.internal.generateMessageAdditions(this, arguments, true);

            if(args !== null) {
            _baseConsole.warn.apply(_baseConsole, args);
            this.backends.send("warn", args);
        }
    },

    error: function() {
        var args = this.internal.generateMessageAdditions(this, arguments, true);

        if(args !== null) {
            _baseConsole.error.apply(_baseConsole, args);
            this.backends.send("error", args);
        }
    },

    // Backends 

    backends: {
        _backends: {},

        add: function(name, backendObject) {
            this._backends[name] = backendObject;
        },

        remove: function(name) {
            delete this._backends[name];
        },

        send: function(function_name, args) {
            for (backend in this._backends) {
                this._backends[backend].send(function_name, args);
            }
        }
    },

    // Stacks 

    stacks: {
        _stacks: [],
        display: false,

        get: function(index) {
            return this._stacks[index];
        }
    },

    // Tags 

    tags: {
        _tags: [],
        _selectedTags: [],
        _muteAllByDefault:true,
        _author: undefined,
        _authorFoldersDisplay: 2,

        add: function(tag){
            this._tags.push(tag);
        },

        remove: function(tag){
            for (var i = 0, l = this._tags.length; i < l; i++) {
                if(this._tags[i] === tag) {
                    this._tags.splice(i, 1);
                    return;
                }
            }
        },

        select: function(tag) {
            this._selectedTags.push(tag);
        },

        unselect: function(tag) {
            for (var i = 0, l = this._selectedTags.length; i < l; i++) {
                if(this._selectedTags[i] === tag) {
                    this._selectedTags.splice(i, 1);
                    return;
                }
            }
        },

        flush: function(){
            this._tags = [];
        },
    },

    // Message filtering 

    filterMessages: function(regex) {
        this._filter = regex;
    },

    // Settings Management 

    settings: {
        _savedProperties: ["tags._selectedTags", "style._colors", "stacks.display","tags._authorFoldersDisplay", "tags._muteAllByDefault"],

        save: function() {
            for (var i = 0, l = this._savedProperties.length; i < l; i++) {
                var currentPropStr = "console." + this._savedProperties[i];

                var val = eval(currentPropStr);
                if(typeof val === "string") {
                    val = "'" + val +"'";
                } else if(typeof val === "object") {
                    val = JSON.stringify(val);
                }

                localStorage.setItem(currentPropStr, val);
            }
        },

        load: function() {
            for (var i = 0, l = this._savedProperties.length; i < l; i++) {
                var currentPropStr = "console." + this._savedProperties[i];

                eval(currentPropStr + "= " + localStorage.getItem(currentPropStr));
            }
        },

        reset: function() {
            for (var i = 0, l = this._savedProperties.length; i < l; i++) {
                var currentPropStr = "console." + this._savedProperties[i];

                eval(currentPropStr + "= " + undefined);
            }
        }
    },

    style: {
        _colors: false
    }
};
console.init();

