import { computed, onMounted, ref, unref, watch } from 'vue';

import { patternsFieldsToForm } from '@/helpers/entities/pattern/fields/form';
import { usePendingHandler } from '@/hooks/query/pending';
import { usePatternsFieldsStore } from '@/stores/modules/patterns-fields';

/**
 * Pattern fields store facade.
 *
 * @returns {Object}
 */
export const usePatternsFields = () => {
  const store = usePatternsFieldsStore();

  return {
    fetchFlappingRulePatternFields: store.fetchFlappingRulePatternFields,
    fetchIdleRulePatternFields: store.fetchIdleRulePatternFields,
    fetchLinkRulePatternFields: store.fetchLinkRulePatternFields,
    fetchResolveRulePatternFields: store.fetchResolveRulePatternFields,
    fetchPbehaviorPatternFields: store.fetchPbehaviorPatternFields,
    fetchAlarmTagPatternFields: store.fetchAlarmTagPatternFields,
    fetchWidgetFilterPatternFields: store.fetchWidgetFilterPatternFields,
    fetchServicePatternFields: store.fetchServicePatternFields,
    fetchStateSettingPatternFields: store.fetchStateSettingPatternFields,
    fetchEventFilterPatternFields: store.fetchEventFilterPatternFields,
    fetchScenarioPatternFields: store.fetchScenarioPatternFields,
    fetchMetaalarmrulePatternFields: store.fetchMetaalarmrulePatternFields,
    fetchDeclareTicketRulePatternFields: store.fetchDeclareTicketRulePatternFields,
    fetchInstructionPatternFields: store.fetchInstructionPatternFields,
    fetchKpiFilterPatternFields: store.fetchKpiFilterPatternFields,
    fetchDynamicInfosPatternFields: store.fetchDynamicInfosPatternFields,
    fetchEventRecordPatternFields: store.fetchEventRecordPatternFields,
  };
};

/**
 * Fetches pattern fields and exposes form-ready attributes.
 *
 * @param {Function} fetchAction
 * @param {boolean} readonly
 * @returns {Object}
 */
export const usePatternsFieldsFetching = (fetchAction, readonly = false) => {
  const patternsFields = ref({});

  const { pending, handler: fetchPatternsFields } = usePendingHandler(async () => {
    patternsFields.value = patternsFieldsToForm(await unref(fetchAction)());
  });

  const alarmAttributes = computed(() => patternsFields.value.alarm_pattern ?? []);
  const entityAttributes = computed(() => patternsFields.value.entity_pattern ?? []);
  const eventAttributes = computed(() => patternsFields.value.event_pattern ?? []);
  const pbehaviorAttributes = computed(() => patternsFields.value.pbehavior_pattern ?? []);
  const weatherServiceAttributes = computed(
    () => patternsFields.value.weather_service_pattern ?? [],
  );

  watch(
    () => unref(fetchAction),
    () => fetchPatternsFields(),
  );

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
