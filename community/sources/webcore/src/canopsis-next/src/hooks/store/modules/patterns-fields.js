import {
  ref,
  unref,
  computed,
  watch,
  onMounted,
} from 'vue';

import { patternsFieldsToForm } from '@/helpers/entities/pattern/fields/form';

import { useStoreModuleHooks } from '@/hooks/store';
import { usePendingHandler } from '@/hooks/query/pending';

/**
 * Creates hooks for accessing the patterns-fields Vuex store module.
 * Provides access to actions for fetching pattern fields for different rule types.
 *
 * @returns {Object} An object containing store module utilities:
 * @property {import('vuex').Store} store - The Vuex store instance
 * @property {import('vuex').Module} module - The patterns-fields module instance
 * @property {Function} useGetters - Function to access module getters
 * @property {Function} useActions - Function to access module actions
 */
const usePatternsFieldsStoreModule = () => useStoreModuleHooks('patternsFields');

/**
 * Custom hook for pattern fields operations.
 * Provides convenient access to pattern fields actions for various rule types.
 *
 * @returns {Object} An object containing pattern fields actions:
 * @property {Function} fetchFlappingRulePatternFields - Fetches pattern fields for flapping rules
 * @property {Function} fetchIdleRulePatternFields - Fetches pattern fields for idle rules
 * @property {Function} fetchLinkRulePatternFields - Fetches pattern fields for link rules
 * @property {Function} fetchResolveRulePatternFields - Fetches pattern fields for resolve rules
 * @property {Function} fetchPbehaviorPatternFields - Fetches pattern fields for pbehaviors
 * @property {Function} fetchAlarmTagPatternFields - Fetches pattern fields for alarm tags
 * @property {Function} fetchWidgetFilterPatternFields - Fetches pattern fields for widget filters
 * @property {Function} fetchServicePatternFields - Fetches pattern fields for entity services
 * @property {Function} fetchStateSettingPatternFields - Fetches pattern fields for state settings
 * @property {Function} fetchEventFilterPatternFields - Fetches pattern fields for event filters
 * @property {Function} fetchScenarioPatternFields - Fetches pattern fields for scenarios
 * @property {Function} fetchMetaalarmrulePatternFields - Fetches pattern fields for meta alarm rules
 * @property {Function} fetchDeclareTicketRulePatternFields - Fetches pattern fields for declare ticket rules
 * @property {Function} fetchInstructionPatternFields - Fetches pattern fields for instructions
 * @property {Function} fetchKpiFilterPatternFields - Fetches pattern fields for KPI filters
 * @property {Function} fetchDynamicInfosPatternFields - Fetches pattern fields for dynamic infos
 * @property {Function} fetchEventRecordPatternFields - Fetches pattern fields for event records
 *
 * @example
 * // Usage in a component
 * const { fetchFlappingRulePatternFields } = usePatternsFields();
 * const fields = await fetchFlappingRulePatternFields({ params: {} });
 */
export const usePatternsFields = () => {
  const { useActions } = usePatternsFieldsStoreModule();

  const actions = useActions({
    fetchFlappingRulePatternFields: 'fetchFlappingRulePatternFields',
    fetchIdleRulePatternFields: 'fetchIdleRulePatternFields',
    fetchLinkRulePatternFields: 'fetchLinkRulePatternFields',
    fetchResolveRulePatternFields: 'fetchResolveRulePatternFields',
    fetchPbehaviorPatternFields: 'fetchPbehaviorPatternFields',
    fetchAlarmTagPatternFields: 'fetchAlarmTagPatternFields',
    fetchWidgetFilterPatternFields: 'fetchWidgetFilterPatternFields',
    fetchServicePatternFields: 'fetchServicePatternFields',
    fetchStateSettingPatternFields: 'fetchStateSettingPatternFields',
    fetchEventFilterPatternFields: 'fetchEventFilterPatternFields',
    fetchScenarioPatternFields: 'fetchScenarioPatternFields',
    fetchMetaalarmrulePatternFields: 'fetchMetaalarmrulePatternFields',
    fetchDeclareTicketRulePatternFields: 'fetchDeclareTicketRulePatternFields',
    fetchInstructionPatternFields: 'fetchInstructionPatternFields',
    fetchKpiFilterPatternFields: 'fetchKpiFilterPatternFields',
    fetchDynamicInfosPatternFields: 'fetchDynamicInfosPatternFields',
    fetchEventRecordPatternFields: 'fetchEventRecordPatternFields',
  });

  return {
    ...actions,
  };
};

/**
 * Hook for fetching and managing pattern fields with automatic loading on mount.
 * Transforms fetched pattern fields data into form-ready format and provides
 * computed properties for alarm and entity pattern fields.
 *
 * @param {Function} [fetchAction = () => ({})] - Async function that fetches pattern fields data.
 *                                               Should return a promise resolving to pattern fields object.
 * @returns {Object} An object containing pattern fields state and computed properties
 *
 * @example
 * // Usage in a component
 * import { usePatternsFields } from '@/hooks/store/modules/patterns-fields';
 *
 * export default {
 *   setup() {
 *     const { fetchFlappingRulePatternFields } = usePatternsFields();
 *     const { pending, alarmAttributes, entityAttributes } = usePatternsFieldsFetching(
 *       fetchFlappingRulePatternFields
 *     );
 *
 *     return {
 *       pending,
 *       alarmPatternFields,
 *       entityPatternFields,
 *     };
 *   },
 * };
 */
export const usePatternsFieldsFetching = (fetchAction, readonly = false) => {
  const patternsFields = ref({});

  const {
    pending,
    handler: fetchPatternsFields,
  } = usePendingHandler(async () => {
    patternsFields.value = patternsFieldsToForm(await unref(fetchAction)());
  });

  const alarmAttributes = computed(() => patternsFields.value.alarm_pattern ?? []);
  const entityAttributes = computed(() => patternsFields.value.entity_pattern ?? []);
  const eventAttributes = computed(() => patternsFields.value.event_pattern ?? []);
  const pbehaviorAttributes = computed(() => patternsFields.value.pbehavior_pattern ?? []);
  const weatherServiceAttributes = computed(() => patternsFields.value.weather_service_pattern ?? []);

  watch(() => unref(fetchAction), () => fetchPatternsFields());

  onMounted(() => {
    if (!readonly) {
      fetchPatternsFields();
    }
  });

  return {
    pending,
    patternsFields,
    alarmAttributes,
    entityAttributes,
    eventAttributes,
    pbehaviorAttributes,
    weatherServiceAttributes,
    fetchPatternsFields,
  };
};
