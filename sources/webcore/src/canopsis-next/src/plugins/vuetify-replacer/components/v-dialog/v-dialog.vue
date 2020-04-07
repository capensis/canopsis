<script>
import VDialog from 'vuetify/es5/components/VDialog';

import overlayableMixin from '../../mixins/overlayable';

export default {
  extends: VDialog,
  mixins: [overlayableMixin],
};
</script>
