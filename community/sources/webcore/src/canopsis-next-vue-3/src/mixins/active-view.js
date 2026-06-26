import { mapState, mapStores } from 'pinia';

import { useActiveViewStore } from '@/stores/modules/active-view';

export const activeViewMixin = {
  computed: {
    ...mapStores(useActiveViewStore),
    ...mapState(useActiveViewStore, {
      view: 'item',
      activeViewPending: 'pending',
    }),
  },
};
