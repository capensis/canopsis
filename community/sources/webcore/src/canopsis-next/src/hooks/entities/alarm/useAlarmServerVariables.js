import { computed } from 'vue';

import { ALARM_PAYLOADS_VARIABLES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useInfosServerVariables } from '@/hooks/entities/common/useInfosServerVariables';

/**
 * Provides a reactive list of alarm server variables with their corresponding localization texts.
 *
 * @returns {Object}
 */
export const useAlarmServerVariables = ({ infos } = {}) => {
  const { t, tc } = useI18n();
  const { variables: infosVariables } = useInfosServerVariables(infos);

  const variables = computed(() => [
    {
      value: ALARM_PAYLOADS_VARIABLES.id,
      text: t('common.alarmId'),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.displayName,
      text: t('alarm.fields.displayName'),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.stateValue,
      text: t('common.state'),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.statusValue,
      text: t('common.status'),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.connector,
      text: t('common.connector'),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.component,
      text: t('common.component'),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.connectorName,
      text: t('common.connectorName'),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.resource,
      text: t('common.resource'),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.infos,
      text: t('common.infos'),
      variables: infosVariables.value,
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.initialOutput,
      text: t('alarm.fields.initialOutput'),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.output,
      text: t('alarm.fields.output'),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.parents,
      text: tc('alarm.fields.parent'),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.children,
      text: tc('alarm.fields.children'),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.totalStateChanges,
      text: tc('alarm.fields.totalStateChanges'),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.eventsCount,
      text: tc('alarm.fields.eventsCount'),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.tags,
      text: tc('common.tag', 2),
    },
    {
      value: ALARM_PAYLOADS_VARIABLES.impactState,
      text: t('common.impactState'),
    },
  ]);

  return {
    variables,
  };
};
