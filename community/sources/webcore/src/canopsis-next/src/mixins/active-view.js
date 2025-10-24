import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('activeView');

/**
 * @mixin Helpers for the view tab entity
 */
export const activeViewMixin = {
  computed: {
    ...mapGetters({
      view: 'item',
      pending: 'pending',
      editing: 'editing',
      editingProcess: 'editingProcess',
    }),
  },
  methods: {
    ...mapActions({
      toggleEditing: 'toggleEditing',
      registerEditingOffHandler: 'registerEditingOffHandler',
      unregisterEditingOffHandler: 'unregisterEditingOffHandler',
      fetchActiveView: 'fetch',
      clearActiveView: 'clear',
    }),
  },
};
