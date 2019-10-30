import { find } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('pbehavior/comment');

/**
 * @mixin
 */
export default {
  methods: {
    ...mapActions({
      createPbehaviorComment: 'create',
      updatePbehaviorComment: 'update',
      removePbehaviorComment: 'remove',
    }),

    updateSeveralPbehaviorComments({ pbehavior, comments }) {
      const oldComments = pbehavior.comments || [];

      /**
       * We are finding comments for creation (without _id field)
       */
      const requestsForCreation = comments.filter(comment => !comment._id)
        .map(comment => this.createPbehaviorComment({
          pbehaviorId: pbehavior._id,
          data: comment,
        }));

      /**
       * We are finding changed comments for updating (with _id field and with changes)
       */
      const requestsForUpdating = comments.filter((comment) => {
        const oldComment = find(oldComments, { _id: comment._id });

        return oldComment && oldComment.message !== comment.message;
      }).map(comment => this.updatePbehaviorComment({
        id: comment._id,
        pbehaviorId: pbehavior._id,
        data: comment,
      }));

      /**
       * We are finding removed comments for removing (with _id)
       */
      const requestsForRemoving = oldComments.filter(oldComment => !find(comments, { _id: oldComment._id }))
        .map(comment => this.removePbehaviorComment({
          id: comment._id,
          pbehaviorId: pbehavior._id,
        }));

      return Promise.all([...requestsForCreation, ...requestsForUpdating, ...requestsForRemoving]);
    },
  },
};
