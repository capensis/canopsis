<script>
import VListGroup from 'vuetify/es5/components/VList/VListGroup';

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
