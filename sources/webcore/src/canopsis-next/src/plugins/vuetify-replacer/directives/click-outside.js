import ClickOutside from 'vuetify/es5/directives/click-outside';
import { getZIndex } from 'vuetify/es5/util/helpers';

import { get } from 'lodash';

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
  const handler = typeof binding.value === 'function' ? binding.value : get(binding, 'value.handler');
  const include = binding.args ? binding.args.include : binding.value.include;
  const { closeConditional } = binding.args ? binding.args : binding.value;

  const defaultCond = get(binding, 'modifiers.zIndex') ? defaultConditionalWithZIndex : defaultConditional;

  const isActive = closeConditional || defaultCond;

  if (!e || isActive(e, el) === false) return;

  if (('isTrusted' in e && !e.isTrusted) || ('pointerType' in e && !e.pointerType)) {
    return;
  }

  const elements = (include || (() => []))();

  elements.push(el);

  if (!elements.some(item => item.contains(e.target))) {
    setTimeout(() => isActive(e, el) && handler && handler(e), 0);
  }
}

export default {
  inserted(el, binding) {
    const onClick = e => directive(e, el, binding);
    const app = document.querySelector('[data-app]') || document.body;
    app.addEventListener('click', onClick, true);
    // eslint-disable-next-line
    el._clickOutside = onClick;
  },

  unbind: ClickOutside.unbind,
};
