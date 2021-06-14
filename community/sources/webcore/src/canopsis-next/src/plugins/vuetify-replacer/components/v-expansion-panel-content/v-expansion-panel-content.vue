<script>
import VExpansionPanelContent from 'vuetify/es5/components/VExpansionPanel/VExpansionPanelContent';

import lazyWithUnmountMixin from '../../mixins/lazy-with-unmount';

export default {
  extends: VExpansionPanelContent,
  mixins: [lazyWithUnmountMixin],
};
</script>
