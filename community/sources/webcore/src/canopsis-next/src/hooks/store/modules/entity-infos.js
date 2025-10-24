import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates an instance of the entity infos store module hooks.
 *
 * @returns {Object} Store module hooks for the 'infos' namespace
 */
export const useEntityInfosStoreModule = () => useStoreModuleHooks('infos');

/**
 * Hook for accessing and managing entity information.
 * Provides access to alarm infos, rules, and entity information.
 *
 * @typedef {Object} EntityInfosGetters
 * @property {Ref<Object>} alarmInfos - Alarm information data
 * @property {Ref<Object>} alarmInfosRules - Rules for alarm information
 * @property {Ref<Object>} entityInfos - Entity information data
 * @property {Ref<boolean>} infosPending - Loading state of entity infos
 *
 * @typedef {Object} EntityInfosActions
 * @property {Function} fetchInfos - Fetches entity information
 *
 * @returns {Object} Hook return object
 * @property {EntityInfosGetters} getters - All available getters
 * @property {EntityInfosActions} actions - All available actions
 */
export const useEntityInfos = () => {
  const { useGetters, useActions } = useEntityInfosStoreModule();

  const getters = useGetters({
    alarmInfos: 'alarmInfos',
    alarmInfosRules: 'alarmInfosRules',
    entityInfos: 'entityInfos',
    infosPending: 'pending',
  });

  const actions = useActions({
    fetchInfos: 'fetch',
  });

  return {
    ...getters,
    ...actions,
  };
};
