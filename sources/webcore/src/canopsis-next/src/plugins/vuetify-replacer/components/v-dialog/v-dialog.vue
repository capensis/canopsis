<script>
import VDialog from 'vuetify/es5/components/VDialog';
import { getZIndex } from 'vuetify/es5/util/helpers';

import overlayableMixin from '../../mixins/overlayable';

export default {
  extends: VDialog,
  mixins: [overlayableMixin],
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
  },
  methods: {
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
};
</script>
