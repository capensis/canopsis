import { createNamespacedHelpers } from 'vuex';

const { mapGetters: sideBarMapGetters } = createNamespacedHelpers('sideBar');
const { mapGetters: modalMapGetters } = createNamespacedHelpers('modals');

export const sideBarInnerMixin = {
  computed: {
    ...sideBarMapGetters({
      sideBarName: 'name',
      sideBarConfig: 'config',
      isSideBarHidden: 'hidden',
    }),

    ...modalMapGetters(['hasMaximizedModal']),
  },
};
