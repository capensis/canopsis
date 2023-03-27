import { getZIndex } from 'vuetify/es5/util/helpers';

import { get } from 'lodash';

import { MIN_CLICK_OUTSIDE_DELAY_AFTER_REGISTERED } from '@/config';

function defaultConditional() {
  return false;
}

/**
 * Calculate zIndexes and return conditional
 *
 * @param e
 * @param el
 * @return {boolean}
 */
function defaultConditionalWithZIndex(e, el) {
  const targetZIndex = getZIndex(e.target);
  const contentZIndex = getZIndex(el);

  return targetZIndex < contentZIndex;
}

/**
 * Function for support old logic v-click-outside and new version. We can use now this directive.
 * With zIndex modification. We can use directive as 'v-click-outside.zIndex'.
 *
 * @param e
 * @param el
 * @param binding
 */
function directive(e, el, binding) {
  // eslint-disable-next-line no-underscore-dangle
  const delayAfterRegistered = Date.now() - el._outsideRegistredAt;

  if (delayAfterRegistered < MIN_CLICK_OUTSIDE_DELAY_AFTER_REGISTERED) {
    return;
  }

  const handler = typeof binding.value === 'function' ? binding.value : get(binding, 'value.handler');
  const include = binding.args ? binding.args.include : binding.value.include;
  const { closeConditional } = binding.args ? binding.args : binding.value;

  const defaultCond = get(binding, 'modifiers.zIndex') ? defaultConditionalWithZIndex : defaultConditional;

  const isActive = closeConditional || defaultCond;

  if (!e || isActive(e, el) === false) return;

  if (('isTrusted' in e && !e.isTrusted) || ('pointerType' in e && !e.pointerType)) {
    return;
  }

  const elements = include ? include() : [];

  elements.push(el);

  if (!elements.some(item => item.contains(e.target))) {
    setTimeout(() => isActive(e, el) && handler && handler(e), 0);
  }
}

/* eslint-disable no-underscore-dangle, no-param-reassign */
export default {
  inserted(el, binding) {
    const onClick = e => directive(e, el, binding);
    const app = document.querySelector('[data-app]') || document.body;
    const { same, zIndex, contextmenu } = binding.modifiers;

    if (same) {
      let mousedownWasOnElement = false;

      const mousedownOutside = (e) => {
        mousedownWasOnElement = el.contains(e.target) || (zIndex && !defaultConditionalWithZIndex(e.target, el));
      };

      const mouseupOutside = (e) => {
        if (!mousedownWasOnElement) {
          onClick(e);
        }
      };

      app.addEventListener('mousedown', mousedownOutside, true);
      app.addEventListener('mouseup', mouseupOutside, true);

      el._mousedownOutside = mousedownOutside;
      el._mouseupOutside = mouseupOutside;
    } else {
      app.addEventListener(contextmenu ? 'contextmenu' : 'click', onClick, true);

      el._clickOutside = onClick;
    }

    el._outsideRegistredAt = Date.now();
  },

  unbind(el, binding) {
    const app = document.querySelector('[data-app]') || document.body; // This is only for unit tests

    if (!app) {
      return;
    }

    const { same, contextmenu } = binding.modifiers;

    if (same) {
      if (el._mousedownOutside) {
        app.removeEventListener('mousedown', el._mousedownOutside, true);
        delete el._mousedownOutside;
      }

      if (el._mouseupOutside) {
        app.removeEventListener('mouseup', el._mouseupOutside, true);
        delete el._mouseupOutside;
      }
    } else if (el._clickOutside) {
      app.removeEventListener(contextmenu ? 'contextmenu' : 'click', el._clickOutside, true);
      delete el._clickOutside;
    }
  },
};
/* eslint-enable no-underscore-dangle, no-param-reassign */
