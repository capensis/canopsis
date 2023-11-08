import { get } from 'lodash';
import { getZIndex } from 'vuetify/lib/util/helpers';
import { attachedRoot } from 'vuetify/lib/util/dom';

import { MIN_CLICK_OUTSIDE_DELAY_AFTER_REGISTERED } from '@/config';

function defaultConditional() {
  return false;
}

function handleShadow(el, callback) {
  const root = attachedRoot(el);
  callback(document);

  if (typeof ShadowRoot !== 'undefined' && root instanceof ShadowRoot) {
    callback(root);
  }
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

function checkIsActive(e, el, binding) {
  const { closeConditional } = binding.args ? binding.args : binding.value;

  const defaultCond = get(binding, 'modifiers.zIndex') ? defaultConditionalWithZIndex : defaultConditional;

  const isActive = closeConditional || defaultCond;

  return isActive(e, el);
}

function checkEvent(e, el, binding) {
  if (!e || checkIsActive(e, el, binding) === false) return false;

  const root = attachedRoot(el);

  if (
    typeof ShadowRoot !== 'undefined'
    && root instanceof ShadowRoot
    && root.host === e.target
  ) {
    return false;
  }

  const { zIndex } = binding.modifiers;

  const include = binding.args
    ? binding.args.include
    : binding.value.include;

  const elements = include ? include() : [];

  elements.push(el);

  const isElementsContainsTarget = elements.some(element => element.contains(e.target));
  const isClickUnder = zIndex ? defaultConditionalWithZIndex(e, el) : true;

  return !isElementsContainsTarget && isClickUnder;
}

/**
 * Function for support old logic v-click-outside and new version. We can use now this directive.
 * With zIndex modification. We can use directive as 'v-click-outside.zIndex'.
 *
 * @param e
 * @param el
 * @param binding
 */
function directive(e, el, binding, vnode) {
  // eslint-disable-next-line no-underscore-dangle
  const clickOutside = el._clickOutside;
  // eslint-disable-next-line no-underscore-dangle
  const nodeClickOutside = clickOutside[vnode.context._uid];

  // eslint-disable-next-line no-underscore-dangle
  const delayAfterRegistered = Date.now() - nodeClickOutside._outsideRegistredAt;

  if (delayAfterRegistered < MIN_CLICK_OUTSIDE_DELAY_AFTER_REGISTERED) {
    return;
  }

  const handler = typeof binding.value === 'function' ? binding.value : get(binding, 'value.handler');

  if (clickOutside.lastMousedownWasOutside && checkEvent(e, el, binding)) {
    setTimeout(() => {
      if (checkIsActive(e, el, binding) && handler) {
        handler(e);
      }
    }, 0);
  }
}

/* eslint-disable no-underscore-dangle, no-param-reassign */
export default {
  inserted(el, binding, vnode) {
    const onClick = e => directive(e, el, binding, vnode);
    const onMousedown = (e) => {
      el._clickOutside.lastMousedownWasOutside = checkEvent(
        e,
        el,
        binding,
      );
    };

    const { same, contextmenu } = binding.modifiers;

    handleShadow(el, (app) => {
      if (same) {
        app.addEventListener('mouseup', onClick, true);
        app.addEventListener('mousedown', onMousedown, true);
      } else {
        app.addEventListener(contextmenu ? 'contextmenu' : 'click', onClick, true);
      }
    });

    if (!el._clickOutside) {
      el._clickOutside = {
        lastMousedownWasOutside: true,
      };
    }

    el._clickOutside[vnode.context._uid] = {
      onClick,
      onMousedown,
      _outsideRegistredAt: Date.now(),
    };
  },

  unbind(el, binding, vnode) {
    const { same, contextmenu } = binding.modifiers;

    handleShadow(el, (app) => {
      const clickOutside = el._clickOutside;

      if (!app || !clickOutside) {
        return;
      }

      if (!app || !clickOutside[vnode.context._uid]) {
        return;
      }

      const {
        onClick,
        onMousedown,
      } = el._clickOutside[vnode.context._uid];

      if (same) {
        app.removeEventListener('mouseup', onClick, true);
        app.removeEventListener('mousedown', onMousedown, true);
      } else if (contextmenu) {
        app.removeEventListener('contextmenu', onClick, true);
      } else {
        app.removeEventListener('click', onClick, true);
      }
    });
  },
};
/* eslint-enable no-underscore-dangle, no-param-reassign */
