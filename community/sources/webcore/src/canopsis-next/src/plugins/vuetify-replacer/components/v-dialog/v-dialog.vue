<script>
import VDialog from 'vuetify/es5/components/VDialog';
import ThemeProvider from 'vuetify/es5/util/ThemeProvider';
import { getZIndex, convertToUnit } from 'vuetify/es5/util/helpers';

import overlayableMixin from '../../mixins/overlayable';

export default {
  extends: VDialog,
  mixins: [overlayableMixin],
  props: {
    customCloseConditional: {
      type: Function,
      default: null,
    },
    ignoreClickOutside: {
      type: Boolean,
      default: false,
    },
    absolute: {
      type: Boolean,
      default: false,
    },
    contentWrapperClass: {
      type: String,
      required: false,
    },
  },
  computed: {
    activeZIndex() {
      if (typeof window === 'undefined') {
        return 0;
      }

      const content = this.stackElement || this.$refs.content;

      /**
       * Return current zindex if not active
       *
       * We've changed factor from 2 to 12 for better overlay animation
       */
      const index = !this.isActive
        ? getZIndex(content)
        : this.getMaxZIndex(this.stackExclude || [content]) + 12;

      if (index == null) {
        return index;
      }

      return parseInt(index, 10);
    },

    contentClasses() {
      const classes = {
        'v-dialog__content': true,
        'v-dialog__content--active': this.isActive,
      };

      if (this.contentWrapperClass) {
        classes[this.contentWrapperClass] = true;
      }

      return classes;
    },
  },
  methods: {
    onFocusin(event) {
      const { target } = event;

      if (
        /**
         * It isn't the document or the dialog body
         */
        ![document, this.$refs.content].includes(target)
        /**
         * It isn't inside the dialog body
         */
        && !this.$refs.content.contains(target)
        /**
         * It isn't inside a dependent element (like a menu)
         */
        && !this.getOpenDependentElements().some(el => el.contains(target))
        /**
         * So we must have focused something outside the dialog and its children
         *
         * We need next tick here for correct zIndex comparison
         */
      ) {
        this.$nextTick(() => {
          /**
           *  We're the topmost dialog
           */
          if (this.activeZIndex >= this.getMaxZIndex()) {
            /**
             * Find and focus the first available element inside the dialog
             */
            const focusable = this.$refs.content.querySelectorAll(
              'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])',
            );

            if (focusable.length) {
              focusable[0].focus();
            }
          }
        });
      }
    },
    closeConditional(e) {
      // If the dialog content contains
      // the click event, or if the
      // dialog is not active
      if (!this.isActive || this.$refs.content.contains(e.target)) {
        return false;
      }
      // If we made it here, the click is outside
      // and is active. If persistent, and the
      // click is on the overlay, animate
      if (this.persistent) {
        if (!this.noClickAnimation && this.overlay === e.target) {
          this.animateClick();
        }

        return false;
      }

      // close dialog if !persistent, clicked outside and we're the topmost dialog.
      // Since this should only be called in a capture event (bottom up), we shouldn't need to stop propagation
      return this.activeZIndex >= this.getMaxZIndex();
    },

    hideScroll() {
      if (this.fullscreen) {
        document.documentElement.classList.add('overflow-y-hidden');
      } else {
        overlayableMixin.methods.hideScroll.call(this);
      }
    },

    getMaxZIndex(exclude = []) {
      const base = this.$el;
      // Start with lowest allowed z-index or z-index of
      // base component's element, whichever is greater
      const zis = [this.stackMinZIndex, getZIndex(base)];
      // Convert the NodeList to an array to
      // prevent an Edge bug with Symbol.iterator
      // https://github.com/vuetifyjs/vuetify/issues/2146
      const activeElements = [
        ...document.getElementsByClassName(this.stackClass),

        /**
         * We've added it here for correct zIndex calculation
         */
        ...document.getElementsByClassName('menuable__content__active v-menu__ignore-click-upper-outside'),
      ];

      // Get z-index for all active dialogs
      for (let index = 0; index < activeElements.length; index += 1) {
        if (!exclude.includes(activeElements[index])) {
          zis.push(getZIndex(activeElements[index]));
        }
      }

      return Math.max(...zis);
    },
  },
  /**
   * We've replaced render method because we've put changes for the `click-outside` directive value function
   *
   * @param h
   * @returns {*}
   */
  render: function render(h) {
    const children = [];
    const directives = [{ name: 'show', value: this.isActive }];

    if (!this.ignoreClickOutside) {
      directives.push({
        name: 'click-outside',
        value: () => {
          /**
           * We can't move call of customCloseConditional into closeConditional because here this method will call
           * more than once. But we need only once call
           */
          if (this.customCloseConditional && this.customCloseConditional() === false) {
            return;
          }

          this.isActive = false;
        },
        args: {
          closeConditional: this.closeConditional,
          include: this.getOpenDependentElements,
        },
        modifiers: { same: true },
      });
    }

    const data = {
      class: this.classes,
      ref: 'dialog',
      directives,
      on: {
        click: (e) => {
          e.stopPropagation();
        },
      },
    };

    if (!this.fullscreen) {
      data.style = {
        maxWidth: this.maxWidth === 'none' ? undefined : convertToUnit(this.maxWidth),
        width: this.width === 'auto' ? undefined : convertToUnit(this.width),
      };
    }

    children.push(this.genActivator());

    let dialog = h('div', data, this.showLazyContent(this.$slots.default));

    if (this.transition) {
      dialog = h('transition', {
        props: {
          name: this.transition,
          origin: this.origin,
        },
      }, [dialog]);
    }

    children.push(h('div', {
      class: this.contentClasses,
      attrs: { tabIndex: '-1', ...this.getScopeIdAttrs() },
      on: {
        keydown: this.onKeydown,
      },
      style: { zIndex: this.activeZIndex },
      ref: 'content',
    }, [this.$createElement(ThemeProvider, {
      props: {
        root: true,
        light: this.light,
        dark: this.dark,
      },
    }, [dialog])]));

    return h('div', {
      staticClass: 'v-dialog__container',
      style: {
        display: !this.hasActivator || this.fullWidth ? 'block' : 'inline-block',
      },
    }, children);
  },
};
</script>
