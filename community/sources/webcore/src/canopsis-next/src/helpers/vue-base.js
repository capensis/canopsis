/**
 * Check if child component has parent component in the $parent
 *
 * @param {VueComponent} child
 * @param {VueComponent} parent
 * @returns {boolean}
 */
export function isParent(child, parent) {
  if (child) {
    if (child === parent || child._original === parent || parent.$el?.contains(child?.$el)) {
      return true;
    }

    if (child.$parent) {
      return isParent(child.$parent, parent);
    }
  }

  return false;
}
