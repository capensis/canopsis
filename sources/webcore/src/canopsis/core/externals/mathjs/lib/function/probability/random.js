'use strict';

module.exports = function (math) {
  var distribution = require('./distribution')(math);

  /**
   * Return a random number between `min` and `max` using a uniform distribution.
   *
   * Syntax:
   *
   *     math.random()                // generate a random number between 0 and 1
   *     math.random(max)             // generate a random number between 0 and max
   *     math.random(min, max)        // generate a random number between min and max
   *     math.random(size)            // generate a matrix with random numbers between 0 and 1
   *     math.random(size, max)       // generate a matrix with random numbers between 0 and max
   *     math.random(size, min, max)  // generate a matrix with random numbers between min and max
   *
   * Examples:
   *
   *     math.random();       // returns a random number between 0 and 1
   *     math.random(100);    // returns a random number between 0 and 100
   *     math.random(30, 40); // returns a random number between 30 and 40
   *     math.random([2, 3]); // returns a 2x3 matrix with random numbers between 0 and 1
   *
   * See also:
   *
   *     randomInt, pickRandom
   *
   * @param {Number} [size] If provided, an array with `size` number of random values is returned
   * @param {Number} [min]  Minimum boundary for the random value
   * @param {Number} [max]  Maximum boundary for the random value
   * @return {Number | Array | Matrix} A random number
   */
  math.random = distribution('uniform').random;
};
