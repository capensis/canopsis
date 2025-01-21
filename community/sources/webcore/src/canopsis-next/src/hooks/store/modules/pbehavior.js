import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the pbehavior Vuex store module.
 * Provides access to pbehavior-specific getters and actions with proper namespacing.
 *
 * @returns {Object} Store module hooks for the 'pbehavior' namespace
 */
const usePbehaviorStoreModule = () => useStoreModuleHooks('pbehavior');

/**
 * Composable hook for accessing pbehavior store functionality.
 * Provides a unified interface to pbehavior-related getters and actions.
 *
 * @returns {Object} Combined pbehavior getters and actions
 *
 * @property {Object} Getters
 * @property {Array<Object>} pbehaviors - List of all pbehavior items
 * @property {Function} getPbehavior - Function to get a specific pbehavior by ID
 * @property {boolean} pbehaviorsPending - Loading state for pbehavior operations
 * @property {Object} pbehaviorsMeta - Metadata for pbehavior list
 *
 * @property {Object} Actions
 * @property {Function} fetchPbehaviorsList - Fetches and stores pbehavior list
 * @property {Function} fetchPbehaviorsListWithoutStore - Fetches pbehaviors without storing
 * @property {Function} fetchPbehaviorsListWithPreviousParams - Refetches with previous parameters
 * @property {Function} fetchPbehaviorEIDSListWithoutStore - Fetches EIDS list without storing
 * @property {Function} createPbehavior - Creates a new pbehavior
 * @property {Function} bulkCreatePbehaviors - Creates multiple pbehaviors
 * @property {Function} createEntityPbehaviors - Creates pbehaviors for entities
 * @property {Function} removeEntityPbehaviors - Removes pbehaviors from entities
 * @property {Function} updatePbehavior - Updates an existing pbehavior
 * @property {Function} bulkUpdatePbehaviors - Updates multiple pbehaviors
 * @property {Function} removePbehavior - Removes a pbehavior
 * @property {Function} bulkRemovePbehaviors - Removes multiple pbehaviors
 * @property {Function} fetchPbehaviorsByEntityId - Fetches pbehaviors for an entity
 * @property {Function} fetchPbehaviorsByEntityIdWithoutStore - Fetches entity pbehaviors without storing
 * @property {Function} fetchPbehaviorsCalendarWithoutStore - Fetches calendar data without storing
 * @property {Function} fetchEntitiesPbehaviorsCalendarWithoutStore - Fetches entities calendar without storing
 */
export const usePbehavior = () => {
  const { useGetters, useActions } = usePbehaviorStoreModule();

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

  return {
    ...getters,
    ...actions,
  };
};
