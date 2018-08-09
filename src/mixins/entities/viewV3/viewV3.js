import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('viewV3');

/**
 * @mixin Helpers for the view entity
 */
export default {
  methods: {
    ...mapActions({
      // fetchView: 'fetchItem',
      createView: 'create',
    }),
  },
  computed: {
    ...mapGetters({
      // view: 'item',
    }),
  },
};
