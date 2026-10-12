import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the dbExport Vuex store module
 * @function useDbExportStoreModule
 * @returns {Object} Store module hooks object
 * @property {import('vuex').Store} store - Vuex store instance
 * @property {import('vuex').Module} module - dbExport module instance
 * @property {Function} useGetters - Hook for accessing module getters
 * @property {Function} useActions - Hook for accessing module actions that make API calls to:
 *   - /api/v4/pbehaviors-db-export
 *   - /api/v4/eventfilter-db-export
 *   - /api/v4/link-rules-db-export
 *   - /api/v4/idle-rules-db-export
 *   - /api/v4/flapping-rules-db-export
 *   - /api/v4/scenarios-db-export
 *   - /api/v4/cat/dynamic-infos-db-export
 *   - /api/v4/cat/declare-ticket-rules-db-export
 *   - /api/v4/cat/metaalarmrules-db-export
 *   - /api/v4/cat/instructions-db-export
 *   - /api/v4/cat/job-configs-db-export
 *   - /api/v4/cat/jobs-db-export
 */
const useDbExportStoreModule = () => useStoreModuleHooks('dbExport');

/**
 * Hook for accessing database export actions from the Vuex store
 * @function useDbExport
 * @returns {Object} Object containing database export actions
 * @property {Function} exportPbehaviorsDb - Exports pbehaviors database
 * @property {Function} exportEventFiltersDb - Exports event filters database
 * @property {Function} exportLinkRulesDb - Exports link rules database
 * @property {Function} exportIdleRulesDb - Exports idle rules database
 * @property {Function} exportFlappingRulesDb - Exports flapping rules database
 * @property {Function} exportScenariosDb - Exports scenarios database
 * @property {Function} exportDynamicInfosDb - Exports dynamic infos database
 * @property {Function} exportDeclareTicketsDb - Exports declare tickets database
 * @property {Function} exportMetaAlarmRulesDb - Exports meta alarm rules database
 * @property {Function} exportInstructionsDb - Exports instructions database
 * @property {Function} exportJobsDb - Exports jobs database
 * @property {Function} exportJobConfigsDb - Exports job configs database
 */
export const useDbExport = () => {
  const { useActions } = useDbExportStoreModule();

  const actions = useActions([
    'exportPbehaviorsDb',
    'exportEventFiltersDb',
    'exportLinkRulesDb',
    'exportIdleRulesDb',
    'exportFlappingRulesDb',
    'exportScenariosDb',
    'exportDynamicInfosDb',
    'exportDeclareTicketsDb',
    'exportMetaAlarmRulesDb',
    'exportInstructionsDb',
    'exportJobsDb',
    'exportJobConfigsDb',
  ]);

  return {
    ...actions,
  };
};
