import getProp from 'lodash/get';

export default function get(object, property, functionFilter) {
  const propertyValue = getProp(object, property);
  if (functionFilter) {
    return functionFilter(property, propertyValue);
  }
  return propertyValue;
}
