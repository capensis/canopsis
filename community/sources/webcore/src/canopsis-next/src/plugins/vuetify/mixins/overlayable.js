import { addPassiveEventListener } from 'vuetify/lib/util/helpers';

export default {
  methods: {
    hideScroll() {
      this.htmlAlreadyContainsOverflowYClass = this.htmlAlreadyContainsOverflowYClass
        ?? document.documentElement.classList.contains('overflow-y-hidden');

      if (this.$vuetify.breakpoint.smAndDown && !this.htmlAlreadyContainsOverflowYClass) {
        document.documentElement.classList.add('overflow-y-hidden');
      } else {
        addPassiveEventListener(window, 'wheel', this.scrollListener, { passive: false });
        window.addEventListener('keydown', this.scrollListener);
      }
    },
    showScroll() {
      if (!this.htmlAlreadyContainsOverflowYClass) {
        document.documentElement.classList.remove('overflow-y-hidden');
      }

      window.removeEventListener('wheel', this.scrollListener);
      window.removeEventListener('keydown', this.scrollListener);

      this.htmlAlreadyContainsOverflowYClass = null;
    },
  },
};
