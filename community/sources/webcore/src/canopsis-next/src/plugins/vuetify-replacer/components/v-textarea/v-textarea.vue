<script>
import VTextarea from 'vuetify/es5/components/VTextarea';

export default {
  extends: VTextarea,
  watch: {
    autoGrow(autoGrow) {
      if (autoGrow) {
        setTimeout(this.calculateInputHeight, 0);
      }
    },
  },
};
</script>
