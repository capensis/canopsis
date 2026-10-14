import { upperFirst } from 'lodash';

import { API_ROUTES } from '@/config';

import request from '@/services/request';

/**
 * Vuex module for database export operations
 * @module dbExport
 *
 * @example
 * // Import and register the module in your Vuex store
 * import dbExport from './store/modules/entities/db-export';
 *
 * const store = new Vuex.Store({
 *   modules: {
 *     dbExport
 *   }
 * });
 *
 * // Usage examples:
 * // Export pbehaviors
 * await store.dispatch('dbExport/exportPbehaviorsDb', { ids: ['123', '456'] });
 *
 * // Export event filters
 * await store.dispatch('dbExport/exportEventFiltersDb');
 *
 * // Export link rules with specific IDs
 * await store.dispatch('dbExport/exportLinkRulesDb', {
 *   ids: ['rule1', 'rule2']
 * });
 *
 * @property {boolean} namespaced - Indicates if module is namespaced
 * @property {Object} actions - Collection of actions for database exports
 * @property {Function} actions.exportPbehaviorsDb - Export pbehaviors
 * @property {Function} actions.exportEventFiltersDb - Export event filters
 * @property {Function} actions.exportLinkRulesDb - Export link rules
 * @property {Function} actions.exportIdleRulesDb - Export idle rules
 * @property {Function} actions.exportFlappingRulesDb - Export flapping rules
 * @property {Function} actions.exportScenariosDb - Export scenarios
 * @property {Function} actions.exportDynamicInfosDb - Export dynamic infos
 * @property {Function} actions.exportDeclareTicketsDb - Export declare tickets
 * @property {Function} actions.exportMetaAlarmRulesDb - Export meta alarm rules
 * @property {Function} actions.exportInstructionsDb - Export instructions
 * @property {Function} actions.exportJobsDb - Export jobs
 * @property {Function} actions.exportJobConfigsDb - Export job configs
 */
export default {
  namespaced: true,
  actions: Object.entries(API_ROUTES.dbExport).reduce((acc, [key, value]) => {
    acc[`export${upperFirst(key)}Db`] = (context, { ids = [] } = {}) => (
      request.post(value, { ids }, { responseType: 'blob' })
    );

    return acc;
  }, {}),
};
