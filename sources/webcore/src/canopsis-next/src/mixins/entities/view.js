import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view');

/**
 * @mixin Helpers for the view entity
 */
export default {
  computed: {
    ...mapGetters({
      view: 'item',
    }),
  },
  methods: {
    ...mapActions({
      fetchView: 'fetchItem',
      createView: 'create',
    }),
  },
};
