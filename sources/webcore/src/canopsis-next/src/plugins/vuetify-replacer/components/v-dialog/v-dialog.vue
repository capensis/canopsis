<script>
import VDialog from 'vuetify/es5/components/VDialog';

import overlayableMixin from '../../mixins/overlayable';

export default {
  extends: VDialog,
  mixins: [overlayableMixin],
  props: {
    onClickOutside: {
      type: Function,
      default: null,
    },
  },
  methods: {
    closeConditional(e) {
      if (!this.isActive || this.$refs.content.contains(e.target)) {
        return false;
      }

      if (this.persistent) {
        if (!this.noClickAnimation && this.overlay === e.target) {
          this.animateClick();
        }

        return false;
      }

      return this.activeZIndex >= this.getMaxZIndex() && this.onClickOutside(e);
    },
  },
};
</script>
