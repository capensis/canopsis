<script>
import VMenu from 'vuetify/lib/components/VMenu';
import Resize from 'vuetify/lib/directives/resize';
import { getZIndex } from 'vuetify/lib/util/helpers';

import ClickOutside from '../../directives/click-outside';

export default {
  directives: {
    ClickOutside,
    Resize,
  },
  extends: VMenu,
  props: {
    ignoreClickUpperOutside: {
      type: Boolean,
      default: false,
    },
    ignoreClickOutside: {
      type: Boolean,
      default: false,
    },
    scrollCalculator: {
      type: Function,
      required: false,
    },
  },
  methods: {
    /**
     * We've added offsetOverflow to condition
     *
     * @param menuWidth
     * @return {string}
     */
    calcLeft(menuWidth) {
      return `${
        this.isAttached && !this.offsetOverflow
          ? this.computedLeft
          : this.calcXOverflow(this.computedLeft, menuWidth)
      }px`;
    },

    calcScrollPosition() {
      const $el = this.$refs.content;
      const activeTile = $el.querySelector('.v-list-item--active');
      const maxScrollTop = $el.scrollHeight - $el.offsetHeight;

      if (activeTile) {
        const newScrollTop = (activeTile.offsetTop - ($el.offsetHeight / 2)) + (activeTile.offsetHeight / 2);

        return Math.min(maxScrollTop, Math.max(0, newScrollTop));
      }

      return this.scrollCalculator ? this.scrollCalculator($el) : $el.scrollTop;
    },

    closeConditional(e) {
      const targetZIndex = getZIndex(e.target);
      const contentZIndex = getZIndex(this.$refs.content);

      if (this.ignoreClickOutside) {
        return false;
      }

      return this.ignoreClickUpperOutside
        ? targetZIndex < contentZIndex
        : this.isActive
        // eslint-disable-next-line no-underscore-dangle
          && !this._isDestroyed
          && this.closeOnClick
          && !this.$refs.content.contains(e.target);
    },

    genContent() {
      const options = {
        attrs: {
          ...this.getScopeIdAttrs(),
          ...this.contentProps,
          role: 'role' in this.$attrs ? this.$attrs.role : 'menu',
        },
        staticClass: 'v-menu__content',
        class: {
          ...this.rootThemeClasses,
          ...this.roundedClasses,
          'v-menu__content--auto': this.auto,
          'v-menu__content--fixed': this.activatorFixed,
          'v-menu__ignore-click-upper-outside': this.ignoreClickUpperOutside,
          menuable__content__active: this.isActive,
          [this.contentClass.trim()]: true,
        },
        style: this.styles,
        directives: this.genDirectives(),
        ref: 'content',
        on: {
          click: (e) => {
            e.stopPropagation();
            if (e.target.getAttribute('disabled')) return;
            if (this.closeOnContentClick) this.isActive = false;
          },
          keydown: this.onKeyDown,
        },
      };

      if (this.$listeners.scroll) {
        options.on = options.on || {};
        options.on.scroll = this.$listeners.scroll;
      }

      if (!this.disabled && this.openOnHover) {
        options.on = options.on || {};
        options.on.mouseenter = this.mouseEnterHandler;
      }

      if (this.openOnHover) {
        options.on = options.on || {};
        options.on.mouseleave = this.mouseLeaveHandler;
      }

      return this.$createElement('div', options, this.getContentSlot());
    },
  },
};
</script>
