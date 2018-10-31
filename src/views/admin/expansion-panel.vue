<script>
import VExpansionPanelContent from 'vuetify/es5/components/VExpansionPanel/VExpansionPanelContent';

export default {
  extends: VExpansionPanelContent,
  watch: {
    isActive(value) {
      if (!value) {
        setTimeout(() => {
          if (!this.isActive) {
            this.isBooted = false;
          }
        }, this.$config.VUETIFY_ANIMATION_DELAY);
      }
    },
  },
};
</script>
