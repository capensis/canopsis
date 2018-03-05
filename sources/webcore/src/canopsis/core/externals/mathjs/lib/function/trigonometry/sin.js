'use strict';

module.exports = function (math) {
  var util = require('../../util/index'),

      BigNumber = math.type.BigNumber,
      Complex = require('../../type/Complex'),
      Unit = require('../../type/Unit'),
      collection = require('../../type/collection'),

      isNumber = util.number.isNumber,
      isBoolean = util['boolean'].isBoolean,
      isComplex = Complex.isComplex,
      isUnit = Unit.isUnit,
      isCollection = collection.isCollection;

  /**
   * Calculate the sine of a value.
   *
   * For matrices, the function is evaluated element wise.
   *
   * Syntax:
   *
   *    math.sin(x)
   *
   * Examples:
   *
   *    math.sin(2);                      // returns Number 0.9092974268256813
   *    math.sin(math.pi / 4);            // returns Number 0.7071067811865475
   *    math.sin(math.unit(90, 'deg'));   // returns Number 1
   *    math.sin(math.unit(30, 'deg'));   // returns Number 0.5
   *
   *    var angle = 0.2;
   *    math.pow(math.sin(angle), 2) + math.pow(math.cos(angle), 2); // returns Number ~1
   *
   * See also:
   *
   *    cos, tan
   *
   * @param {Number | Boolean | Complex | Unit | Array | Matrix | null} x  Function input
   * @return {Number | Complex | Array | Matrix} Sine of x
   */
  math.sin = function sin(x) {
    if (arguments.length != 1) {
      throw new math.error.ArgumentsError('sin', arguments.length, 1);
    }

    if (isNumber(x)) {
      return Math.sin(x);
    }

    if (isComplex(x)) {
      return new Complex(
          0.5 * Math.sin(x.re) * (Math.exp(-x.im) + Math.exp( x.im)),
          0.5 * Math.cos(x.re) * (Math.exp( x.im) - Math.exp(-x.im))
      );
    }

    if (isUnit(x)) {
      if (!x.hasBase(Unit.BASE_UNITS.ANGLE)) {
        throw new TypeError ('Unit in function sin is no angle');
      }
      return Math.sin(x.value);
    }

    if (isCollection(x)) {
      return collection.deepMap(x, sin);
    }

    if (isBoolean(x) || x === null) {
      return Math.sin(x);
    }

    if (x instanceof BigNumber) {
      // TODO: implement BigNumber support
      // downgrade to Number
      return sin(x.toNumber());
    }

    throw new math.error.UnsupportedTypeError('sin', math['typeof'](x));
  };
};
