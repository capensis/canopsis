import { useStoreModuleHooks } from '@/hooks/store';

const useEntityCommentStoreModule = () => useStoreModuleHooks('entity/comment');

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
