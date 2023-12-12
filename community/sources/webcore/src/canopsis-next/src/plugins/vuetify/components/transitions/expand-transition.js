import { upperFirst } from 'vuetify/es5/util/helpers';

/* eslint-disable no-underscore-dangle, no-param-reassign */
const defineProperty = (obj, key, value) => {
  if (key in obj) {
    Object.defineProperty(obj, key, { value, enumerable: true, configurable: true, writable: true });
  } else {
    obj[key] = value;
  }

  return obj;
};

export default (expandedParentClass = '', x = false) => {
  const sizeProperty = x ? 'width' : 'height';

  function resetStyles(el) {
    el.style.overflow = el._initialStyle.overflow;
    el.style[sizeProperty] = el._initialStyle[sizeProperty];
    delete el._initialStyle;
  }

  function afterLeave(el) {
    if (expandedParentClass && el._parent) {
      el._parent.classList.remove(expandedParentClass);
    }

    resetStyles(el);
  }

  return {
    beforeEnter: function beforeEnter(el) {
      el._initialStyle = defineProperty({
        transition: el.style.transition,
        visibility: el.style.visibility,
        overflow: el.style.overflow,
      }, sizeProperty, el.style[sizeProperty]);
    },
    enter: function enter(el) {
      const initialStyle = el._initialStyle;

      el._parent = el.parentNode;
      el.style.setProperty('transition', 'none', 'important');
      el.style.visibility = 'hidden';

      const size = `${el[`offset${upperFirst(sizeProperty)}`]}px`;

      el.style.visibility = initialStyle.visibility;
      el.style.overflow = 'hidden';
      el.style[sizeProperty] = 0;
      el.offsetHeight; // force reflow
      el.style.transition = initialStyle.transition;

      if (expandedParentClass && el._parent) {
        el._parent.classList.add(expandedParentClass);
      }

      requestAnimationFrame(() => {
        el.style[sizeProperty] = size;
      });
    },

    afterEnter: resetStyles,
    enterCancelled: resetStyles,
    leave: function leave(el) {
      el._initialStyle = defineProperty({
        overflow: el.style.overflow,
      }, sizeProperty, el.style[sizeProperty]);
      el.style.overflow = 'hidden';
      el.style[sizeProperty] = `${el[`offset${upperFirst(sizeProperty)}`]}px`;
      el.offsetHeight; // force reflow
      requestAnimationFrame(() => el.style[sizeProperty] = 0);
    },

    afterLeave,
    leaveCancelled: afterLeave,
  };
};
/* eslint-enable no-underscore-dangle, no-param-reassign */
