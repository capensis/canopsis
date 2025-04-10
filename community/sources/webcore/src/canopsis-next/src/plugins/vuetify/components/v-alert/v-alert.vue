<script>
import VAlert from 'vuetify/lib/components/VAlert/VAlert';

export default {
  extends: VAlert,
  computed: {
    computedColor() {
      return this.color || `var(--v-${this.type}-background-base)`;
    },

    iconColor() {
      return `var(--v-${this.type}-base)`;
    },
  },
};
</script>

<style lang="scss">
.theme--dark.v-sheet.v-alert {
  color: #212121;
  font-weight: 700;
  font-size: 14px;
  line-height: 22px;
}
</style>
