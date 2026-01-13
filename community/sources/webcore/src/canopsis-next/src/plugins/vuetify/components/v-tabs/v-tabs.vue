<script>
import { throttle } from 'lodash';
import VTabs from 'vuetify/lib/components/VTabs/VTabs';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

export default {
  extends: VTabs,
  watch: {
    '$i18n.locale': {
      handler() {
        this.$nextTick(this.onResize);
      },
    },
  },
  created() {
    this.onResize = throttle(this.onResizeOriginal, VUETIFY_ANIMATION_DELAY);
  },
  methods: {
    onResizeOriginal() {
      // eslint-disable-next-line no-underscore-dangle
      if (this._isDestroyed) return;
      clearTimeout(this.resizeTimeout);
      this.resizeTimeout = window.setTimeout(this.callSlider, 0);
    },
  },
};
</script>
