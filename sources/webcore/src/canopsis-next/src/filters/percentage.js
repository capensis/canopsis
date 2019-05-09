export default function (value, precision = 3) {
  const tmp = 10 ** precision;

  const filteredValue = 100 * value;

  const roundedValue = Math.round(filteredValue * tmp) / tmp;

  return `${roundedValue}%`;
}
