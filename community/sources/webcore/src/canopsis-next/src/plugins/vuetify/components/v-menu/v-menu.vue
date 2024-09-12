<script>
import VMenu from 'vuetify/es5/components/VMenu';
import { getZIndex } from 'vuetify/es5/util/helpers';

import lazyWithUnmountMixin from '../../mixins/lazy-with-unmount';

export default {
  extends: VMenu,
  mixins: [lazyWithUnmountMixin],
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
      const activeTile = $el.querySelector('.v-list__tile--active');
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
          && this.closeOnClick
          && !this.$refs.content.contains(e.target);
    },

    genContent() {
      const options = {
        attrs: this.getScopeIdAttrs(),
        staticClass: 'v-menu__content',
        class: {
          ...this.rootThemeClasses,
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

      if (!this.disabled && this.openOnHover) {
        options.on.mouseenter = this.mouseEnterHandler;
      }

      if (this.openOnHover) {
        options.on.mouseleave = this.mouseLeaveHandler;
      }

      return this.$createElement('div', options, this.showLazyContent(this.$slots.default));
    },
  },
};
</script>
