<script>
import VMenu from 'vuetify/es5/components/VMenu';
import { getZIndex } from 'vuetify/es5/util/helpers';

export default {
  extends: VMenu,
  props: {
    ignoreClickUpperOutside: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    closeConditional(e) {
      const targetZIndex = getZIndex(e.target);
      const contentZIndex = getZIndex(this.$refs.content);

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
