import { find } from 'lodash';
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
      fetchPbehaviorsListWithPreviousParams: 'fetchListWithPreviousParams',
      fetchPbehaviorEIDSListWithoutStore: 'fetchEIDSWithoutStore',
      createPbehavior: 'create',
      bulkCreatePbehaviors: 'bulkCreate',
      createEntityPbehaviors: 'bulkCreateEntityPbehaviors',
      removeEntityPbehaviors: 'bulkRemoveEntityPbehaviors',
      updatePbehavior: 'update',
      bulkUpdatePbehaviors: 'bulkUpdate',
      removePbehavior: 'remove',
      bulkRemovePbehaviors: 'bulkRemove',
      fetchPbehaviorsByEntityId: 'fetchListByEntityId',
      fetchPbehaviorsByEntityIdWithoutStore: 'fetchListByEntityIdWithoutStore',
      fetchPbehaviorsCalendarWithoutStore: 'fetchPbehaviorsCalendarWithoutStore',
      fetchEntitiesPbehaviorsCalendarWithoutStore: 'fetchEntitiesPbehaviorsCalendarWithoutStore',
    }),

    async createPbehaviorWithComments({ data }) {
      const pbehavior = await this.createPbehavior({ data });

      await this.updateSeveralPbehaviorComments({ comments: data.comments, pbehavior });

      return pbehavior;
    },

    async createPbehaviorsWithComments(pbehaviors) {
      const response = await this.bulkCreatePbehaviors({ data: pbehaviors });

      await Promise.all(
        response.map(({ id, errors, item: pbehavior }) => {
          if (!errors) {
            return this.updateSeveralPbehaviorComments({
              comments: pbehavior.comments,
              pbehavior: {
                ...pbehavior,
                _id: id,
                comments: [],
              },
            });
          }

          return Promise.reject(errors);
        }),
      );

      return response;
    },

    async updatePbehaviorsWithComments(pbehaviors = [], originalPbehaviors = []) {
      const response = await this.bulkUpdatePbehaviors({ data: pbehaviors });

      await Promise.all(
        pbehaviors.map(pbehavior => this.updateSeveralPbehaviorComments({
          pbehavior: find(originalPbehaviors, { _id: pbehavior._id }),
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
