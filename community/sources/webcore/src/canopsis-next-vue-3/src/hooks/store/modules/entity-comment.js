import { useEntityCommentStore } from '@/stores/modules/entity-comment';

/**
 * Entity comment store facade.
 *
 * @returns {Object}
 */
export const useEntityComment = () => {
  const entityCommentStore = useEntityCommentStore();

  return {
    createEntityComment: entityCommentStore.create,
    updateEntityComment: entityCommentStore.update,
    fetchEntityCommentsListWithoutStore: entityCommentStore.fetchListWithoutStore,
  };
};
