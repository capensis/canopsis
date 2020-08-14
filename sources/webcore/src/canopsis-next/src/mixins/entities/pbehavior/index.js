import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('pbehavior');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      pbehaviors: 'items',
      pbehaviorsPending: 'pending',
      pbehaviorsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchPbehaviorsList: 'fetchList',
      fetchPbehaviorsListWithoutStore: 'fetchListWithoutStore',
      createPbehavior: 'create',
      updatePbehavior: 'update',
      removePbehavior: 'remove',
      fetchPbehaviorsByEntityId: 'fetchListByEntityId',
    }),

    createPbehaviors(pbehaviors, options = {}) {
      return Promise.all(pbehaviors.map(data => this.createPbehavior({ data, ...options })));
    },

    removePbehaviors(ids) {
      return Promise.all(ids.map(id => this.removePbehavior({ id })));
    },

    updatePbehaviors(pbehaviors) {
      return Promise.all(pbehaviors.map(data => this.updatePbehavior({ data, id: data._id })));
    },
  },
};
