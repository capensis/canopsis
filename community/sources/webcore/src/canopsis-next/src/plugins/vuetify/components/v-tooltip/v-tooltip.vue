<script>
import VTooltip from 'vuetify/lib/components/VTooltip';
import { omit } from 'lodash';

export default {
  extends: VTooltip,
  props: {
    ignoreContentLeave: {
      type: Boolean,
      default: false,
    },
    customActivator: {
      type: Boolean,
      default: false,
    },
    disableResize: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    computedTransition() {
      if (this.transition) {
        return this.transition;
      }

      if (this.top) {
        return 'slide-y-reverse-transition';
      }

      if (this.right) {
        return 'slide-x-transition';
      }
      if (this.bottom) {
        return 'slide-y-transition';
      }
      if (this.left) {
        return 'slide-x-reverse-transition';
      }
      return '';
    },
  },
  mounted() {
    if (this.disableResize) {
      window?.removeEventListener('resize', this.updateDimensions, false);
    }
  },
  methods: {
    mouseEnterHandlers(e) {
      this.getActivator(e);
      this.runDelay('open');
    },

    genActivatorListeners() {
      let listeners = VTooltip.options.methods.genActivatorListeners.call(this);

      if (this.ignoreContentLeave) {
        listeners = omit(listeners, ['mouseleave']);
      }

      return listeners;
    },
  },
};
</script>
