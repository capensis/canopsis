import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view');

/**
 * @mixin Helpers for the view entity
 */
export default {
  computed: {
    ...mapGetters({
      view: 'item',
      views: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchView: 'fetchItem',
      createView: 'create',
      updateView: 'update',
      removeView: 'remove',
    }),
  },
};
