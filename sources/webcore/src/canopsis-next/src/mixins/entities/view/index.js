import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view');

/**
 * @mixin Helpers for the view entity
 */
export default {
  computed: {
    ...mapGetters({
      viewId: 'itemId',
      viewPending: 'pending',
      view: 'item',
      getViewById: 'getItemById',
    }),
  },
  methods: {
    ...mapActions({
      fetchView: 'fetchItem',
      createView: 'create',
      updateView: 'update',
      updateViewWithoutStore: 'updateWithoutStore',
      removeView: 'remove',
    }),
  },
};
