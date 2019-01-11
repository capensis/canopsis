<script>
import VExpansionPanelContent from 'vuetify/es5/components/VExpansionPanel/VExpansionPanelContent';

export default {
  extends: VExpansionPanelContent,
  props: {
    ...VExpansionPanelContent.props,

    lazyWithUnmount: {
      type: Boolean,
      default: false,
    },
  },
  watch: {
    isActive(value) {
      if (this.lazyWithUnmount && !value) {
        /**
         * Animation waiting and after that unmounting component content
         */
        setTimeout(() => {
          if (!this.isActive) {
            this.isBooted = false;
          }
        }, this.$config.VUETIFY_ANIMATION_DELAY * 2);
      }
    },
  },
};
</script>
