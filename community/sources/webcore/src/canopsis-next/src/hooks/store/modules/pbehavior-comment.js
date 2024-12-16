import { find } from 'lodash';

import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the pbehavior comment Vuex store module
 *
 * @returns {Object} Store module hooks
 * @property {import('vuex').Store} store - Vuex store instance
 * @property {import('vuex').Module} module - pbehavior comment module instance
 * @property {Function} useGetters - Hook for accessing module getters
 * @property {Function} useActions - Hook for accessing module actions
 */
const usePbehaviorCommentStoreModule = () => useStoreModuleHooks('pbehavior/comment');

/**
 * Hook for managing pbehavior comments - creating, updating, removing and bulk operations
 *
 * @returns {Object} Comment management functions
 * @property {Object} actions - Comment store actions
 * @property {Function} actions.createPbehaviorComment - Create single comment
 * @property {Function} actions.updatePbehaviorComment - Update single comment
 * @property {Function} actions.removePbehaviorComment - Remove single comment
 * @property {Function} createComments - Create multiple comments
 * @property {Function} updateComments - Update multiple existing comments
 * @property {Function} removeComments - Remove multiple comments
 * @property {Function} updateSeveralPbehaviorComments - Bulk update comments handling creates/updates/deletes
 */
export const usePbehaviorComment = () => {
  const { t } = useI18n();
  const popups = usePopups();

  const { useActions } = usePbehaviorCommentStoreModule();

  const actions = useActions({
    createPbehaviorComment: 'create',
    updatePbehaviorComment: 'update',
    removePbehaviorComment: 'remove',
  });

  /**
   * Creates multiple comments for a pbehavior
   *
   * @param {Object} params - Function parameters
   * @param {Array<Object>} params.comments - Array of comment objects to create
   * @param {string} params.pbehaviorId - ID of the pbehavior to attach comments to
   * @returns {Array<Promise>} Array of promises for comment creation
   */
  const createComments = ({ comments, pbehaviorId }) => (
    comments.map(comment => actions.createPbehaviorComment({ pbehaviorId, data: comment }))
  );

  /**
   * Updates multiple existing comments
   *
   * @param {Object} params - Function parameters
   * @param {Array<Object>} params.comments - Array of comment objects to update
   * @param {string} params.comments[]._id - ID of comment to update
   * @param {string} params.pbehaviorId - ID of the associated pbehavior
   * @returns {Array<Promise>} Array of promises for comment updates
   */
  const updateComments = ({ comments, pbehaviorId }) => comments.map(comment => actions.updatePbehaviorComment({
    pbehaviorId,
    id: comment._id,
    data: comment,
  }));

  /**
   * Removes multiple comments
   *
   * @param {Object} params - Function parameters
   * @param {Array<Object>} params.comments - Array of comment objects to remove
   * @param {string} params.comments[]._id - ID of comment to remove
   * @returns {Array<Promise>} Array of promises for comment removal
   */
  const removeComments = ({ comments }) => comments.map(comment => actions.removePbehaviorComment({ id: comment._id }));

  /**
   * Updates multiple comments for a pbehavior, handling creates, updates and deletes
   *
   * @param {Object} params - Function parameters
   * @param {Object} params.pbehavior - Pbehavior object containing existing comments
   * @param {string} params.pbehavior._id - ID of the pbehavior
   * @param {Array<Object>} params.pbehavior.comments - Existing comments on the pbehavior
   * @param {Array<Object>} params.comments - New comments array to apply
   * @param {string} [params.comments[]._id] - ID of existing comment (if updating)
   * @param {string} params.comments[].message - Comment message content
   * @throws {Error} If any comment operations fail
   * @returns {Promise<void>}
   */
  const updateSeveralPbehaviorComments = async ({ pbehavior, comments }) => {
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
        ...actions.createComments({ comments: newComments, pbehaviorId: pbehavior._id }),
        ...actions.updateComments({ comments: changedComments, pbehaviorId: pbehavior._id }),
        ...actions.removeComments({ comments: removedComments }),
      ]);
    } catch (err) {
      const message = Object.values(err).join('\n');

      console.error(err);

      popups.error({ text: message || t('errors.default') });
    }
  };

  return {
    ...actions,

    createComments,
    updateComments,
    removeComments,
    updateSeveralPbehaviorComments,
  };
};
