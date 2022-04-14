import { createNamespacedHelpers } from 'vuex';

import entitiesPbehaviorCommentMixin from '@/mixins/entities/pbehavior/comment';

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
      updatePbehavior: 'update',
      removePbehavior: 'remove',
      fetchPbehaviorsByEntityId: 'fetchListByEntityId',
      fetchPbehaviorsByEntityIdWithoutStore: 'fetchListByEntityIdWithoutStore',
      fetchPbehaviorsCalendarWithoutStore: 'fetchPbehaviorsCalendarWithoutStore',
      fetchEntitiesPbehaviorsCalendarWithoutStore: 'fetchEntitiesPbehaviorsCalendarWithoutStore',
    }),

    async createPbehaviors(pbehaviors) {
      return Promise.all(pbehaviors.map(async (data) => {
        const pbehavior = await this.createPbehavior({ data });

        await this.updateSeveralPbehaviorComments({ comments: data.comments, pbehavior });
      }));
    },

    removePbehaviors(pbehaviors) {
      return Promise.all(pbehaviors.map(({ _id }) => this.removePbehavior({ id: _id })));
    },

    updatePbehaviors(pbehaviors) {
      return Promise.all(pbehaviors.map(pbehavior => Promise.all([
        this.updatePbehavior({ data: pbehavior, id: pbehavior._id }),
        this.updateSeveralPbehaviorComments({
          pbehavior: this.getPbehavior(pbehavior._id),
          comments: pbehavior.comments,
        }),
      ])));
    },
  },
};
