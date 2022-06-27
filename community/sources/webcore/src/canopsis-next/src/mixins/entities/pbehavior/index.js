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
      updatePbehavior: 'update',
      removePbehavior: 'remove',
      fetchPbehaviorsByEntityId: 'fetchListByEntityId',
      fetchPbehaviorsByEntityIdWithoutStore: 'fetchListByEntityIdWithoutStore',
    }),

    async createPbehaviorWithComments({ data }) {
      const pbehavior = await this.createPbehavior({ data });

      await this.updateSeveralPbehaviorComments({ comments: data.comments, pbehavior });

      return pbehavior;
    },

    createPbehaviorsWithComments(pbehaviors) {
      return Promise.all(pbehaviors.map(data => this.createPbehaviorWithComments({ data })));
    },

    updatePbehaviorsWithComments(pbehaviors) {
      return Promise.all(pbehaviors.map(data => Promise.all([
        this.updatePbehavior({ id: data._id, data }),
        this.updateSeveralPbehaviorComments({
          pbehavior: this.getPbehavior(data._id),
          comments: data.comments,
        }),
      ])));
    },

    removePbehaviors(pbehaviors) {
      return Promise.all(pbehaviors.map(({ _id }) => this.removePbehavior({ id: _id })));
    },
  },
};
