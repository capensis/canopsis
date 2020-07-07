<script>
import VMenu from 'vuetify/es5/components/VMenu';

export default {
  extends: VMenu,
  props: {
    ignoreClickOutside: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    closeConditional(e) {
      return !this.ignoreClickOutside && this.isActive &&
        this.closeOnClick &&
        !this.$refs.content.contains(e.target);
    },
  },
};
</script>
