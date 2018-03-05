var sc = require( '../' );

sc
    .option( 'foo', {
        abbr: 'f',
        desc: 'description for option foo',
        default: 'bar'
    } )
    .option( 'baz', {
        abbr: 'b',
        desc: 'description for baz flag',
        flag: true
    } );

console.log( sc.parse() );
