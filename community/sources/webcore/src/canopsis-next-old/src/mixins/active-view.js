import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('activeView');

/**
 * @mixin Helpers for the view tab entity
 */
export const activeViewMixin = {
  computed: {
    ...mapGetters({
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
    ...mapActions({
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
