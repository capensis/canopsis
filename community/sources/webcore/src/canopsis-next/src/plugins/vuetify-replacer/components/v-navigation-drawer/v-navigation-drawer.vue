<script>
import VNavigationDrawer from 'vuetify/es5/components/VNavigationDrawer';

import overlayableMixin from '../../mixins/overlayable';

export default {
  extends: VNavigationDrawer,
  mixins: [overlayableMixin],
  props: {
    ignoreClickOutside: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    reactsToClick() {
      return !this.stateless && !this.permanent && (this.isMobile || this.temporary) && !this.ignoreClickOutside;
    },
  },
};
</script>
