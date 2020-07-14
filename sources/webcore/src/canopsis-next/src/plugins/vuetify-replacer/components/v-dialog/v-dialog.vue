<script>
import VDialog from 'vuetify/es5/components/VDialog';
import { getZIndex } from 'vuetify/es5/util/helpers';

import overlayableMixin from '../../mixins/overlayable';

export default {
  extends: VDialog,
  mixins: [overlayableMixin],
  props: {
    zIndex: {
      type: Number,
      default: null,
    },
  },
  computed: {
    activeZIndex() {
      if (typeof window === 'undefined') return 0;

      const content = this.stackElement || this.$refs.content;
      // Return current zindex if not active
      const index = !this.isActive
        ? getZIndex(content)
        : this.zIndex || this.getMaxZIndex(this.stackExclude || [content]) + 2;

      if (index == null) {
        return index;
      }

      // Return max current z-index (excluding self) + 2
      // (2 to leave room for an overlay below, if needed)
      return parseInt(index, 10);
    },
  },

};
</script>
