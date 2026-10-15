/* eslint-disable no-console */

const eachIterator = (arrayOrObject, func) => {
  const array = Array.isArray(arrayOrObject) ? arrayOrObject : Object.entries(arrayOrObject);

  array.forEach(func);
};

const enhanceBenchmarkFunction = (func) => {
  // eslint-disable-next-line no-param-reassign
  func.each = (arrayOrObject, nameFunc, benchmarkFunc) => eachIterator(arrayOrObject, (item) => {
    const title = nameFunc(item);

    func(title, (...args) => benchmarkFunc(item, ...args));
  });
};

module.exports = {
  eachIterator,
  enhanceBenchmarkFunction,
};
