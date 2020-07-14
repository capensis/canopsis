<script>
import VMenu from 'vuetify/es5/components/VMenu';
import { getZIndex } from 'vuetify/es5/util/helpers';

export default {
  extends: VMenu,
  props: {
    ignoreClickUpperOutside: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    closeConditional(e) {
      const targetZIndex = getZIndex(e.target);
      const contentZIndex = getZIndex(this.$refs.content);

      return (this.ignoreClickUpperOutside && targetZIndex < contentZIndex)
        || (!this.ignoreClickUpperOutside
          && this.isActive
          && this.closeOnClick
          && !this.$refs.content.contains(e.target)
        );
    },
  },
};
</script>
