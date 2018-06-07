import getProp from 'lodash/get';

export default function get(object, property) {
  return getProp(object, property);
}
