import { find } from 'lodash';

import { mapIds } from '@/helpers/array';

import { useStoreModuleHooks } from '@/hooks/store';
import { usePbehaviorComment } from '@/hooks/store/modules/pbehavior-comment';

/**
 * Creates hooks for accessing the pbehavior Vuex store module
 *
 * @returns {Object} Store module hooks
 * @property {import('vuex').Store} store - Vuex store instance
 * @property {import('vuex').Module} module - pbehavior module instance
 * @property {Function} useGetters - Hook for accessing module getters
 * @property {Function} useActions - Hook for accessing module actions
 */
const usePbehaviorStoreModule = () => useStoreModuleHooks('pbehavior');

/**
 * Hook for accessing and managing pbehaviors and their comments
 *
 * @returns {Object} Pbehavior management functions and store access
 * @property {Object} getters - Pbehavior store getters
 * @property {Array} getters.pbehaviors - List of pbehavior items
 * @property {Function} getters.getPbehavior - Get single pbehavior by ID
 * @property {boolean} getters.pbehaviorsPending - Loading state
 * @property {Object} getters.pbehaviorsMeta - Pagination metadata
 * @property {Object} actions - Pbehavior store actions
 * @property {Function} actions.fetchPbehaviorsList - Fetch pbehaviors list
 * @property {Function} actions.fetchPbehaviorsListWithoutStore - Fetch without storing
 * @property {Function} actions.createPbehavior - Create single pbehavior
 * @property {Function} actions.bulkCreatePbehaviors - Bulk create pbehaviors
 * @property {Function} actions.updatePbehavior - Update single pbehavior
 * @property {Function} actions.bulkUpdatePbehaviors - Bulk update pbehaviors
 * @property {Function} actions.removePbehavior - Remove single pbehavior
 * @property {Function} actions.bulkRemovePbehaviors - Bulk remove pbehaviors
 * @property {Function} actions.fetchPbehaviorsByEntityId - Get pbehaviors by entity
 * @property {Function} actions.fetchPbehaviorsCalendarWithoutStore - Get calendar data
 */
export const usePbehavior = () => {
  const { useGetters, useActions } = usePbehaviorStoreModule();
  const pbehaviorCommentModule = usePbehaviorComment();

  const getters = useGetters({
    pbehaviors: 'items',
    getPbehavior: 'getItem',
    pbehaviorsPending: 'pending',
    pbehaviorsMeta: 'meta',
  });

  const actions = useActions({
    fetchPbehaviorsList: 'fetchList',
    fetchPbehaviorsListWithoutStore: 'fetchListWithoutStore',
    fetchPbehaviorsListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchPbehaviorEIDSListWithoutStore: 'fetchEIDSWithoutStore',
    createPbehavior: 'create',
    bulkCreatePbehaviors: 'bulkCreate',
    createEntityPbehaviors: 'bulkCreateEntityPbehaviors',
    removeEntityPbehaviors: 'bulkRemoveEntityPbehaviors',
    updatePbehavior: 'update',
    bulkUpdatePbehaviors: 'bulkUpdate',
    removePbehavior: 'remove',
    bulkRemovePbehaviors: 'bulkRemove',
    fetchPbehaviorsByEntityId: 'fetchListByEntityId',
    fetchPbehaviorsByEntityIdWithoutStore: 'fetchListByEntityIdWithoutStore',
    fetchPbehaviorsCalendarWithoutStore: 'fetchPbehaviorsCalendarWithoutStore',
    fetchEntitiesPbehaviorsCalendarWithoutStore: 'fetchEntitiesPbehaviorsCalendarWithoutStore',
  });

  /**
   * Creates a pbehavior with associated comments
   *
   * @param {Object} params - Function parameters
   * @param {Object} params.data - Pbehavior data with comments
   * @param {Object} params.data.comments - Comments to associate with pbehavior
   * @returns {Promise<Object>} Created pbehavior object
   */
  const createPbehaviorWithComments = async ({ data }) => {
    const pbehavior = await actions.createPbehavior({ data });

    await pbehaviorCommentModule.updateSeveralPbehaviorComments({ comments: data.comments, pbehavior });

    return pbehavior;
  };

  /**
   * Bulk creates multiple pbehaviors with their associated comments
   *
   * @param {Array<Object>} pbehaviors - Array of pbehavior objects to create
   * @param {Object} pbehaviors[].comments - Comments to associate with each pbehavior
   * @returns {Promise<Array<Object>>} Array of created pbehaviors with status
   * @property {string} response[].id - ID of created pbehavior
   * @property {Object} response[].item - Created pbehavior data
   * @property {Object} [response[].errors] - Any errors that occurred
   */
  const createPbehaviorsWithComments = async (pbehaviors = []) => {
    const response = await actions.bulkCreatePbehaviors({ data: pbehaviors });

    await Promise.all(
      response.map(({ id, errors, item: pbehavior }) => {
        if (!errors) {
          return pbehaviorCommentModule.updateSeveralPbehaviorComments({
            comments: pbehavior.comments,
            pbehavior: {
              ...pbehavior,
              _id: id,
              comments: [],
            },
          });
        }

        return Promise.reject(errors);
      }),
    );

    return response;
  };

  /**
   * Updates multiple pbehaviors and their associated comments
   *
   * @param {Array<Object>} pbehaviors - Array of pbehavior objects to update
   * @param {string} pbehaviors[]._id - ID of pbehavior to update
   * @param {Array<Object>} pbehaviors[].comments - Updated comments
   * @param {Array<Object>} originalPbehaviors - Original pbehavior objects before update
   * @returns {Promise<Array<Object>>} Array of updated pbehaviors
   */
  const updatePbehaviorsWithComments = async (pbehaviors = [], originalPbehaviors = []) => {
    const response = await actions.bulkUpdatePbehaviors({ data: pbehaviors });

    await Promise.all(
      pbehaviors.map(pbehavior => pbehaviorCommentModule.updateSeveralPbehaviorComments({
        pbehavior: find(originalPbehaviors, { _id: pbehavior._id }),
        comments: pbehavior.comments,
      })),
    );

    return response;
  };

  /**
   * Removes multiple pbehaviors
   *
   * @param {Array<Object>} pbehaviors - Array of pbehavior objects to remove
   * @param {string} pbehaviors[]._id - ID of pbehavior to remove
   * @returns {Promise<void>}
   */
  const removePbehaviors = (pbehaviors = []) => actions.bulkRemovePbehaviors({ data: mapIds(pbehaviors) });

  return {
    ...getters,
    ...actions,
    ...pbehaviorCommentModule,

    createPbehaviorWithComments,
    createPbehaviorsWithComments,
    updatePbehaviorsWithComments,
    removePbehaviors,
  };
};
