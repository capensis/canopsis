import { tourBaseMixin } from './base';

export const tourInnerMixin = {
  mixins: [tourBaseMixin],
  props: {
    callbacks: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    tourCallbacks() {
      return {
        ...this.callbacks,

        onStop: this.onStop,
      };
    },
    tourInstance() {
      return this.$tours[this.tourName];
    },
  },
  mounted() {
    if (this.tourInstance) {
      this.tourInstance.start();
    }
  },
  methods: {
    async onStop() {
      if (this.callbacks.onStop) {
        await this.callbacks.onStop();
      }

      await this.finishTourByName(this.tourName);
    },
  },
};
