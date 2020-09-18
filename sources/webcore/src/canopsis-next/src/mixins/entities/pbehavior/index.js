import { createNamespacedHelpers } from 'vuex';

import entitiesPbehaviorCommentMixin from '@/mixins/entities/pbehavior/comment';

const { mapGetters, mapActions } = createNamespacedHelpers('pbehavior');

/**
 * @mixin
 */
export default {
  mixins: [entitiesPbehaviorCommentMixin],
  computed: {
    ...mapGetters({
      pbehaviors: 'items',
      getPbehavior: 'getItem',
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

    async createPbehaviors(pbehaviors) {
      return Promise.all(pbehaviors.map(async (data) => {
        const pbehavior = await this.createPbehavior({ data });
        return this.updateSeveralPbehaviorComments({ comments: data.comments, pbehavior });
      }));
    },

    removePbehaviors(pbehaviors) {
      return Promise.all(pbehaviors.map(({ _id }) => this.removePbehavior({ id: _id })));
    },

    updatePbehaviors(pbehaviors) {
      return Promise.all(pbehaviors.map((pbehavior) => {
        this.updatePbehavior({ data: pbehavior, id: pbehavior._id });
        return this.updateSeveralPbehaviorComments({
          pbehavior: this.getPbehavior(pbehavior._id),
          comments: pbehavior.comments,
        });
      }));
    },
  },
};
