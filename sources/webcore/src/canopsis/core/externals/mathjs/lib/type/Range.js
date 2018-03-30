'use strict';

var util = require('../util/index');

var number = util.number;
var string = util.string;
var array = util.array;

/**
 * @constructor Range
 * Create a range. A range has a start, step, and end, and contains functions
 * to iterate over the range.
 *
 * A range can be constructed as:
 *     var range = new Range(start, end);
 *     var range = new Range(start, end, step);
 *
 * To get the result of the range:
 *     range.forEach(function (x) {
 *         console.log(x);
 *     });
 *     range.map(function (x) {
 *         return math.sin(x);
 *     });
 *     range.toArray();
 *
 * Example usage:
 *     var c = new Range(2, 6);         // 2:1:5
 *     c.toArray();                     // [2, 3, 4, 5]
 *     var d = new Range(2, -3, -1);    // 2:-1:-2
 *     d.toArray();                     // [2, 1, 0, -1, -2]
 *
 * @param {Number} start  included lower bound
 * @param {Number} end    excluded upper bound
 * @param {Number} [step] step size, default value is 1
 */
function Range(start, end, step) {
  if (!(this instanceof Range)) {
    throw new SyntaxError('Constructor must be called with the new operator');
  }

  if (start != null && !number.isNumber(start)) {
    throw new TypeError('Parameter start must be a number');
  }
  if (end != null && !number.isNumber(end)) {
    throw new TypeError('Parameter end must be a number');
  }
  if (step != null && !number.isNumber(step)) {
    throw new TypeError('Parameter step must be a number');
  }

  this.start = (start != null) ? parseFloat(start) : 0;
  this.end   = (end != null) ? parseFloat(end) : 0;
  this.step  = (step != null) ? parseFloat(step) : 1;
}

/**
 * Parse a string into a range,
 * The string contains the start, optional step, and end, separated by a colon.
 * If the string does not contain a valid range, null is returned.
 * For example str='0:2:11'.
 * @param {String} str
 * @return {Range | null} range
 */
Range.parse = function (str) {
  if (!string.isString(str)) {
    return null;
  }

  var args = str.split(':');
  var nums = args.map(function (arg) {
    return parseFloat(arg);
  });

  var invalid = nums.some(function (num) {
    return isNaN(num);
  });
  if(invalid) {
    return null;
  }

  switch (nums.length) {
    case 2: return new Range(nums[0], nums[1]);
    case 3: return new Range(nums[0], nums[2], nums[1]);
    default: return null;
  }
};

/**
 * Create a clone of the range
 * @return {Range} clone
 */
Range.prototype.clone = function () {
  return new Range(this.start, this.end, this.step);
};

/**
 * Test whether an object is a Range
 * @param {*} object
 * @return {Boolean} isRange
 */
Range.isRange = function (object) {
  return (object instanceof Range);
};

/**
 * Retrieve the size of the range.
 * Returns an array containing one number, the number of elements in the range.
 * @returns {Number[]} size
 */
Range.prototype.size = function () {
  var len = 0,
      start = this.start,
      step = this.step,
      end = this.end,
      diff = end - start;

  if (number.sign(step) == number.sign(diff)) {
    len = Math.ceil((diff) / step);
  }
  else if (diff == 0) {
    len = 0;
  }

  if (isNaN(len)) {
    len = 0;
  }
  return [len];
};

/**
 * Calculate the minimum value in the range
 * @return {Number | undefined} min
 */
Range.prototype.min = function () {
  var size = this.size()[0];

  if (size > 0) {
    if (this.step > 0) {
      // positive step
      return this.start;
    }
    else {
      // negative step
      return this.start + (size - 1) * this.step;
    }
  }
  else {
    return undefined;
  }
};

/**
 * Calculate the maximum value in the range
 * @return {Number | undefined} max
 */
Range.prototype.max = function () {
  var size = this.size()[0];

  if (size > 0) {
    if (this.step > 0) {
      // positive step
      return this.start + (size - 1) * this.step;
    }
    else {
      // negative step
      return this.start;
    }
  }
  else {
    return undefined;
  }
};


/**
 * Execute a callback function for each value in the range.
 * @param {function} callback   The callback method is invoked with three
 *                              parameters: the value of the element, the index
 *                              of the element, and the Matrix being traversed.
 */
Range.prototype.forEach = function (callback) {
  var x = this.start;
  var step = this.step;
  var end = this.end;
  var i = 0;

  if (step > 0) {
    while (x < end) {
      callback(x, i, this);
      x += step;
      i++;
    }
  }
  else if (step < 0) {
    while (x > end) {
      callback(x, i, this);
      x += step;
      i++;
    }
  }
};

/**
 * Execute a callback function for each value in the Range, and return the
 * results as an array
 * @param {function} callback   The callback method is invoked with three
 *                              parameters: the value of the element, the index
 *                              of the element, and the Matrix being traversed.
 * @returns {Array} array
 */
Range.prototype.map = function (callback) {
  var array = [];
  this.forEach(function (value, index, obj) {
    array[index] = callback(value, index, obj);
  });
  return array;
};

/**
 * Create an Array with a copy of the Ranges data
 * @returns {Array} array
 */
Range.prototype.toArray = function () {
  var array = [];
  this.forEach(function (value, index) {
    array[index] = value;
  });
  return array;
};

/**
 * Get the primitive value of the Range, a one dimensional array
 * @returns {Array} array
 */
Range.prototype.valueOf = function () {
  // TODO: implement a caching mechanism for range.valueOf()
  return this.toArray();
};

/**
 * Get a string representation of the range, with optional formatting options.
 * Output is formatted as 'start:step:end', for example '2:6' or '0:0.2:11'
 * @param {Object | Number | Function} [options]  Formatting options. See
 *                                                lib/util/number:format for a
 *                                                description of the available
 *                                                options.
 * @returns {String} str
 */
Range.prototype.format = function (options) {
  var str = number.format(this.start, options);

  if (this.step != 1) {
    str += ':' + number.format(this.step, options);
  }
  str += ':' + number.format(this.end, options);
  return str;
};

/**
 * Get a string representation of the range.
 * @returns {String}
 */
Range.prototype.toString = function () {
  return this.format();
};

// exports
module.exports = Range;
