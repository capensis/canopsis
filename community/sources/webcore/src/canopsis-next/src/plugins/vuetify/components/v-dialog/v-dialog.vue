<script>
import VDialog from 'vuetify/lib/components/VDialog';
import { getZIndex, convertToUnit } from 'vuetify/lib/util/helpers';

import ClickOutside from '../../directives/click-outside';
import overlayableMixin from '../../mixins/overlayable';

export default {
  directives: {
    ClickOutside,
  },
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
      if (!event || !this.retainFocus) return;

      const { target } = event;

      if (
        !!target
        && this.$refs.dialog
        /**
         * It isn't the document or the dialog body
         */
        && ![document, this.$refs.dialog].includes(target)
        /**
         * It isn't inside the dialog body
         */
        && !this.$refs.dialog.contains(target)
        /*
         * It isn't inside a dependent element (like a menu)
         */
        && this.activeZIndex >= this.getMaxZIndex()
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
        /**
         * Find and focus the first available element inside the dialog
         */
        const focusable = this.$refs.dialog.querySelectorAll(
          'button, [href], input:not([type="hidden"]), select, textarea, [tabindex]:not([tabindex="-1"])',
        );

        const focusableElement = [...focusable].find(element => !element.hasAttribute('disabled') && !element.matches('[tabindex="-1"]'));

        if (focusableElement) {
          focusableElement.focus();
        }
      }
    },
    closeConditional(e) {
      const { target } = e; // Ignore the click if the dialog is closed or destroyed,
      // if it was on an element inside the content,
      // if it was dragged onto the overlay (#6969),
      // or if this isn't the topmost dialog (#9907)

      // eslint-disable-next-line no-underscore-dangle
      if (this._isDestroyed || !this.isActive || this.$refs.content.contains(target)) {
        return false;
      }

      if (this.overlay && target && !this.overlay.$el.contains(target)) {
        return false;
      }

      // If we made it here, the click is outside
      // and is active. If persistent, and the
      // click is on the overlay, animate
      if (this.persistent) {
        if (!this.noClickAnimation) {
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
        ...document.getElementsByClassName('v-menu__content--active'),
        ...document.getElementsByClassName('v-dialog__content--active'),

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

    genInnerContent() {
      const directives = [{
        name: 'show',
        value: this.isActive,
      }];

      if (!this.ignoreClickOutside) {
        directives.push({
          name: 'click-outside',
          value: {
            handler: (event) => {
              /**
               * We can't move call of customCloseConditional into closeConditional because here this method will call
               * more than once. But we need only once call
               */
              if (this.customCloseConditional && this.customCloseConditional() === false) {
                return;
              }

              this.onClickOutside(event);
            },
            closeConditional: this.closeConditional,
            include: this.getOpenDependentElements,
          },
          modifiers: { same: true },
        });
      }

      const data = {
        class: this.classes,
        attrs: {
          tabindex: this.isActive ? 0 : undefined,
        },
        ref: 'dialog',
        directives,
        style: {
          transformOrigin: this.origin,
        },
      };

      if (!this.fullscreen) {
        data.style = {
          ...data.style,
          maxWidth: convertToUnit(this.maxWidth),
          width: convertToUnit(this.width),
        };
      }

      return this.$createElement('div', data, this.getContentSlot());
    },
  },

  render(h) {
    return h('div', {
      staticClass: 'v-dialog__container',
      class: {
        'v-dialog__container--attached': this.attach === ''
          || this.attach === true
          || this.attach === 'attach',
      },
    }, [this.genActivator(), this.genContent()]);
  },
};
</script>
