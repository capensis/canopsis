<script>
import VTabs from 'vuetify/es5/components/VTabs/VTabs';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

export default {
  extends: VTabs,
  mounted() {
    setTimeout(this.callSlider, VUETIFY_ANIMATION_DELAY);
  },
};
</script>
