import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view');

/**
 * @mixin Helpers for the view entity
 */
export default {
  computed: {
    ...mapGetters({
      view: 'item',
      getItemById: 'getItemById',
    }),
  },
  methods: {
    ...mapActions({
      fetchView: 'fetchItem',
      createView: 'create',
      updateView: 'update',
      removeView: 'remove',
    }),

    getViewById(id) {
      return this.getItemById(id);
    },
  },
};
