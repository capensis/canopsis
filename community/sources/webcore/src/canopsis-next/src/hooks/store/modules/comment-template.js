import { useStoreModuleHooks } from '@/hooks/store';

const useCommentTemplateStoreModule = () => useStoreModuleHooks('commentTemplate');

/**
 * Custom hook for accessing actions related to comment templates from Vuex store module.
 *
 * @returns {Object} An object containing functions to create, update, remove, and fetch comment templates list.
 */
export const useCommentTemplates = () => {
  const { useActions } = useCommentTemplateStoreModule();

  const actions = useActions({
    createCommentTemplate: 'create',
    updateCommentTemplate: 'update',
    removeCommentTemplate: 'remove',
    fetchCommentTemplatesListWithoutStore: 'fetchListWithoutStore',
  });

  return {
    ...actions,
  };
};
