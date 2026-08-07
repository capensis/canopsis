import { mapActions, mapState, mapStores } from 'pinia';

import { useActiveViewStore } from '@/stores/modules/active-view';

export const activeViewMixin = {
  computed: {
    ...mapStores(useActiveViewStore),
    ...mapState(useActiveViewStore, {
      view: 'item',
      activeViewPending: 'pending',
      activeViewScreenMode: 'screenMode',
      activeViewEditing: 'editing',
      activeViewEditingProcess: 'editingProcess',
      activeViewPeriodicRefreshPaused: 'periodicRefreshPaused',
      activeViewIsKioskScreenMode: 'isKioskScreenMode',
    }),
  },
  methods: {
    ...mapActions(useActiveViewStore, {
      toggleEditing: 'toggleEditing',
      registerEditingOffHandler: 'registerEditingOffHandler',
      unregisterEditingOffHandler: 'unregisterEditingOffHandler',
      fetchActiveView: 'fetch',
      clearActiveView: 'clear',
      resumePeriodicRefresh: 'resumePeriodicRefresh',
      pausePeriodicRefresh: 'pausePeriodicRefresh',
      setActiveViewScreenMode: 'setScreenMode',
    }),
  },
};
