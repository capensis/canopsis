import { storeToRefs } from 'pinia';

import { usePbehaviorStore } from '@/stores/modules/pbehavior';

import { pickIds } from '@/helpers/array';

import { usePbehaviorComment } from '@/hooks/store/modules/pbehavior-comment';

/**
 * Pbehavior store facade.
 *
 * @returns {Object}
 */
export const usePbehavior = () => {
  const pbehaviorStore = usePbehaviorStore();
  const pbehaviorCommentModule = usePbehaviorComment();
  const { items, meta, pending, fetchingParams } = storeToRefs(pbehaviorStore);

  const removePbehaviors = (pbehaviors = []) =>
    pbehaviorStore.bulkRemove({ data: pickIds(pbehaviors) });

  return {
    pbehaviors: items,
    pbehaviorsMeta: meta,
    pbehaviorsPending: pending,
    pbehaviorsFetchingParams: fetchingParams,

    fetchPbehaviorsList: pbehaviorStore.fetchList,
    fetchPbehaviorsListWithPreviousParams: pbehaviorStore.fetchListWithPreviousParams,
    fetchPbehaviorsListWithoutStore: pbehaviorStore.fetchListWithoutStore,
    fetchPbehaviorWithoutStore: pbehaviorStore.fetchItemWithoutStore,
    fetchPbehaviorEIDSListWithoutStore: pbehaviorStore.fetchListWithoutStore,
    createPbehavior: pbehaviorStore.create,
    bulkCreatePbehaviors: pbehaviorStore.bulkCreate,
    createEntityPbehaviors: pbehaviorStore.bulkCreateEntityPbehaviors,
    removeEntityPbehaviors: pbehaviorStore.bulkRemoveEntityPbehaviors,
    updatePbehavior: pbehaviorStore.update,
    bulkUpdatePbehaviors: pbehaviorStore.bulkUpdate,
    bulkEnablePbehaviors: pbehaviorStore.bulkEnable,
    bulkDisablePbehaviors: pbehaviorStore.bulkDisable,
    removePbehavior: pbehaviorStore.remove,
    bulkRemovePbehaviors: pbehaviorStore.bulkRemove,
    fetchPbehaviorsByEntityIdWithoutStore: pbehaviorStore.fetchPbehaviorsByEntityIdWithoutStore,
    fetchPbehaviorsCalendarWithoutStore: pbehaviorStore.fetchPbehaviorsCalendarWithoutStore,
    fetchEntitiesPbehaviorsCalendarWithoutStore:
      pbehaviorStore.fetchEntitiesPbehaviorsCalendarWithoutStore,
    fetchPbehaviorReasonsListWithoutStore: pbehaviorStore.fetchPbehaviorReasonsListWithoutStore,
    fetchPbehaviorTypesListWithoutStore: pbehaviorStore.fetchPbehaviorTypesListWithoutStore,
    fetchPbehaviorTypesFieldList: pbehaviorStore.fetchPbehaviorTypesFieldList,
    fetchNextPbehaviorTypePriority: pbehaviorStore.fetchNextPbehaviorTypePriority,

    ...pbehaviorCommentModule,

    removePbehaviors,
  };
};
