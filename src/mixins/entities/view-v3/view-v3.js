import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('viewV3');

/**
 * @mixin Helpers for the view entity
 */
export default {
  computed: {
    ...mapGetters({
      // view: 'item',
    }),
  },
  methods: {
    ...mapActions({
      // fetchView: 'fetchItem',
      createView: 'create',
    }),
  },
};
