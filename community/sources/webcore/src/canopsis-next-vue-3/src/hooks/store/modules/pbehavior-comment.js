import { usePbehaviorCommentStore } from '@/stores/modules/pbehavior-comment';

/**
 * Pbehavior comment store facade.
 *
 * @returns {Object}
 */
export const usePbehaviorComment = () => {
  const pbehaviorCommentStore = usePbehaviorCommentStore();

  const updateSeveralPbehaviorComments = ({ comments = [], pbehavior } = {}) =>
    Promise.all(
      comments.map((comment) =>
        comment._id
          ? pbehaviorCommentStore.update({
              data: comment,
              pbehaviorId: pbehavior._id,
              commentId: comment._id,
            })
          : pbehaviorCommentStore.create({ data: comment, pbehaviorId: pbehavior._id }),
      ),
    );

  const removeComments = ({ comments = [] } = {}) =>
    comments.map((comment) =>
      pbehaviorCommentStore.remove({
        id: comment._id,
      }),
    );

  return {
    createPbehaviorComment: pbehaviorCommentStore.create,
    updatePbehaviorComment: pbehaviorCommentStore.update,
    removePbehaviorComment: pbehaviorCommentStore.remove,
    updateSeveralPbehaviorComments,
    removeComments,
  };
};
