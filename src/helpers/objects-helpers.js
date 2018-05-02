/**
 * Get the descendant of attribute's object with a string
 * Yon can now access a descendant attribute like that : Object['att1.att2']
 * @param item The Object in which you want the attribute
 * @param stringProp The string that describe the attribute : ['att1.att2']
 * @returns {T} The value corresponding
 */
function getDescendantProp(item, stringProp) {
  return stringProp.split('.').reduce((a, b) => a[b], item);
}

export default { getDescendantProp };
