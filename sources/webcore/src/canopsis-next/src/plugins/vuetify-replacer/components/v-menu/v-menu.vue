<script>
import VMenu from 'vuetify/es5/components/VMenu';

import lazyWithUnmountMixin from '../../mixins/lazy-with-unmount';

export default {
  extends: VMenu,
  mixins: [lazyWithUnmountMixin],
};
</script>
