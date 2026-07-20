import { storeToRefs } from 'pinia'

import { useCommentTemplateStore } from '@/stores/modules/comment-template'

/**
 * Comment template store facade.
 *
 * @returns {Object}
 */
export const useCommentTemplate = () => {
  const commentTemplateStore = useCommentTemplateStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(commentTemplateStore)

  return {
    commentTemplates: items,
    commentTemplatesMeta: meta,
    commentTemplatesPending: pending,
    commentTemplatesFetchingParams: fetchingParams,

    fetchCommentTemplatesList: commentTemplateStore.fetchList,
    fetchCommentTemplatesListWithPreviousParams: commentTemplateStore.fetchListWithPreviousParams,
    fetchCommentTemplatesListWithoutStore: commentTemplateStore.fetchListWithoutStore,
    createCommentTemplate: commentTemplateStore.create,
    updateCommentTemplate: commentTemplateStore.update,
    removeCommentTemplate: commentTemplateStore.remove,
  }
}
