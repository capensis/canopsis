<script>
import VAvatar from 'vuetify/lib/components/VAvatar/VAvatar';

export default {
  extends: VAvatar,
  props: {
    outlined: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    classes() {
      return {
        'v-avatar--left': this.left,
        'v-avatar--right': this.right,
        'v-avatar--outlined': this.outlined,
        ...this.roundedClasses,
      };
    },
  },
};
</script>

<style lang="scss">
.v-avatar--outlined {
  background-color: transparent;
}
</style>
