var sc = require( '../' );

sc
    .command( 'version', {
        desc: 'display app\'s version',
        callback: function () {
            console.log( 'version' );
        }
    } )
    .end()
    .command( 'server', {
        desc: 'handle the server'
    } )
    .option( 'port', {
        abbr: 'p',
        desc: 'Server port',
        default: '8080'
    } )
    .option( 'hostname', {
        abbr: 'H',
        desc: 'Server hostname'
    } )
    .command( 'start', {
        desc: 'start the server',
        callback: function ( options ) {
            var port = options.port,
                hostname = options.hostname;

            console.log( port, hostname );
        }
    } )
    .end()
    .command( 'stop', {
        desc: 'stop the server',
        callback: function () {}
    } );

sc.parse();
