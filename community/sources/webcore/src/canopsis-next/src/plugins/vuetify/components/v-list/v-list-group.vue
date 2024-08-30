<script>
import VListGroup from 'vuetify/lib/components/VList/VListGroup';

export default {
  extends: VListGroup,
  props: {
    hideOnDisabled: {
      type: Boolean,
      default: false,
    },
  },
  watch: {
    disabled(value, oldValue) {
      if (this.hideOnDisabled && value && !oldValue) {
        this.isActive = false;
      }
    },
  },
};
</script>
