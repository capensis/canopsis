import { find } from 'lodash';

import { usePbehaviorCommentStore } from '@/stores/modules/pbehavior-comment';

/**
 * Pbehavior comment store facade.
 *
 * @returns {Object}
 */
export const usePbehaviorComment = () => {
  const pbehaviorCommentStore = usePbehaviorCommentStore();

  const createComments = ({ comments = [], pbehaviorId } = {}) =>
    comments.map((comment) => pbehaviorCommentStore.create({ data: comment, pbehaviorId }));

  const updateComments = ({ comments = [], pbehaviorId } = {}) =>
    comments.map((comment) =>
      pbehaviorCommentStore.update({
        data: comment,
        pbehaviorId,
        commentId: comment._id,
      }),
    );

  const removeComments = ({ comments = [] } = {}) =>
    comments.map((comment) =>
      pbehaviorCommentStore.remove({
        id: comment._id,
      }),
    );

  const updateSeveralPbehaviorComments = async ({ comments = [], pbehavior } = {}) => {
    const oldComments = pbehavior?.comments ?? [];
    const newComments = comments.filter((comment) => !comment._id);
    const changedComments = comments.filter((comment) => {
      const oldComment = find(oldComments, { _id: comment._id });

      return oldComment && oldComment.message !== comment.message;
    });
    const removedComments = oldComments.filter(
      (oldComment) => !find(comments, { _id: oldComment._id }),
    );

    return Promise.all([
      ...createComments({ comments: newComments, pbehaviorId: pbehavior._id }),
      ...updateComments({ comments: changedComments, pbehaviorId: pbehavior._id }),
      ...removeComments({ comments: removedComments }),
    ]);
  };

  return {
    createPbehaviorComment: pbehaviorCommentStore.create,
    updatePbehaviorComment: pbehaviorCommentStore.update,
    removePbehaviorComment: pbehaviorCommentStore.remove,
    createComments,
    updateComments,
    updateSeveralPbehaviorComments,
    removeComments,
  };
};
