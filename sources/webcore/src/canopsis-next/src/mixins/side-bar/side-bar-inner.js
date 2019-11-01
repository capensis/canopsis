import { createNamespacedHelpers } from 'vuex';

import sideBarMixin from './side-bar';

const { mapGetters: sideBarMapGetters } = createNamespacedHelpers('sideBar');
const { mapGetters: modalMapGetters } = createNamespacedHelpers('modals');

export default {
  mixins: [sideBarMixin],
  computed: {
    ...sideBarMapGetters({
      sideBarName: 'name',
      sideBarConfig: 'config',
      isSideBarHidden: 'hidden',
    }),
    ...modalMapGetters(['hasModals']),
  },
};
