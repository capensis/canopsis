import { find } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('pbehavior/comment');

/**
 * @mixin
 */
export const entitiesPbehaviorCommentMixin = {
  methods: {
    ...mapActions({
      createPbehaviorComment: 'create',
      updatePbehaviorComment: 'update',
      removePbehaviorComment: 'remove',
    }),

    createComments({ comments, pbehaviorId }) {
      return comments.map(comment => this.createPbehaviorComment({ pbehaviorId, data: comment }));
    },

    updateComments({ comments, pbehaviorId }) {
      return comments.map(comment => this.updatePbehaviorComment({
        pbehaviorId,
        id: comment._id,
        data: comment,
      }));
    },

    removeComments({ comments }) {
      return comments.map(comment => this.removePbehaviorComment({ id: comment._id }));
    },

    async updateSeveralPbehaviorComments({ pbehavior, comments }) {
      const oldComments = pbehavior?.comments ?? [];

      /**
       * We are finding comments for creation (without _id field)
       */
      const newComments = comments.filter(comment => !comment._id);

      /**
       * We are finding changed comments for updating (with _id field and with changes)
       */
      const changedComments = comments.filter((comment) => {
        const oldComment = find(oldComments, { _id: comment._id });

        return oldComment && oldComment.message !== comment.message;
      });

      /**
       * We are finding removed comments for removing (with _id)
       */
      const removedComments = oldComments.filter(oldComment => !find(comments, { _id: oldComment._id }));

      try {
        await Promise.all([
          ...this.createComments({ comments: newComments, pbehaviorId: pbehavior._id }),
          ...this.updateComments({ comments: changedComments, pbehaviorId: pbehavior._id }),
          ...this.removeComments({ comments: removedComments }),
        ]);
      } catch (err) {
        const message = Object.values(err).join('\n');

        console.error(err);

        this.$popups.error({ text: message || this.$t('errors.default') });
      }
    },
  },
};
