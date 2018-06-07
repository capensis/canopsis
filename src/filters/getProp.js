import getProp from 'lodash/get';

export default function get(object, property, filter) {
  const value = getProp(object, property);
  if (filter) {
    return filter(value);
  }
  return value;
}
