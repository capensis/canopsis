import { useStoreModuleHooks } from '@/hooks/store';

const useEntityCommentStoreModule = () => useStoreModuleHooks('entity/comment');

/**
 * A custom hook for interacting with entity comments in the Vuex store.
 *
 * @returns {Object} An object containing actions to create, update, and fetch entity comments.
 *
 * @example
 * // Usage example
 * import { useEntityComments } from './path/to/useEntityComments';
 *
 * export default {
 *   setup() {
 *     const { createEntityComment, updateEntityComment, fetchEntityCommentsListWithoutStore } = useEntityComments();
 *
 *     // Call the actions as needed
 *     createEntityComment(commentData);
 *     updateEntityComment(commentId, updatedData);
 *     fetchEntityCommentsListWithoutStore(entityId);
 *   }
 * }
 */
export const useEntityComments = () => {
  const { useActions } = useEntityCommentStoreModule();

  const actions = useActions({
    createEntityComment: 'create',
    updateEntityComment: 'update',
    fetchEntityCommentsListWithoutStore: 'fetchListWithoutStore',
  });

  return {
    ...actions,
  };
};
