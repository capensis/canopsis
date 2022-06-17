import { createNamespacedHelpers } from 'vuex';

import { entitiesPbehaviorCommentMixin } from '@/mixins/entities/pbehavior/comment';

const { mapGetters, mapActions } = createNamespacedHelpers('pbehavior');

export const entitiesPbehaviorMixin = {
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
      fetchPbehaviorEIDSListWithoutStore: 'fetchEIDSWithoutStore',
      createPbehavior: 'create',
      bulkCreatePbehaviors: 'bulkCreate',
      updatePbehavior: 'update',
      bulkUpdatePbehaviors: 'bulkUpdate',
      removePbehavior: 'remove',
      bulkRemovePbehaviors: 'bulkRemove',
      fetchPbehaviorsByEntityId: 'fetchListByEntityId',
      fetchPbehaviorsByEntityIdWithoutStore: 'fetchListByEntityIdWithoutStore',
    }),

    async createPbehaviorWithComments({ data }) {
      const pbehavior = await this.createPbehavior({ data });

      await this.updateSeveralPbehaviorComments({ comments: data.comments, pbehavior });

      return pbehavior;
    },

    async createPbehaviorsWithComments(pbehaviors) {
      const response = await this.bulkCreatePbehaviors({ data: pbehaviors });

      await Promise.all(
        response.map(({ id, item: pbehavior }) => this.updateSeveralPbehaviorComments({
          comments: pbehavior.comments,
          pbehavior: {
            ...pbehavior,
            _id: id,
            comments: [],
          },
        })),
      );

      return response;
    },

    async updatePbehaviorsWithComments(pbehaviors) {
      const response = await this.bulkUpdatePbehaviors({ data: pbehaviors });

      await Promise.all(
        pbehaviors.map(pbehavior => this.updateSeveralPbehaviorComments({
          pbehavior: this.getPbehavior(pbehavior._id),
          comments: pbehavior.comments,
        })),
      );

      return response;
    },

    removePbehaviors(pbehaviors) {
      return this.bulkRemovePbehaviors({
        data: pbehaviors.map(({ _id }) => ({ _id })),
      });
    },
  },
};
