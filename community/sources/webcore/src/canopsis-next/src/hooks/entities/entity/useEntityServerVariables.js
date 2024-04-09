import { computed } from 'vue';

import { ENTITY_PAYLOADS_VARIABLES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useInfosServerVariables } from '@/hooks/entities/common/useInfosServerVariables';

/**
 * Provides a reactive list of entity server variables with their corresponding localization texts.
 *
 * @returns {Object}
 */
export const useEntityServerVariables = ({ infos } = {}) => {
  const { t, tc } = useI18n();

  const { variables: infosVariables } = useInfosServerVariables(infos);
  const variables = computed(() => [
    {
      value: ENTITY_PAYLOADS_VARIABLES.id,
      text: t('common.entityId'),
    },
    {
      value: ENTITY_PAYLOADS_VARIABLES.name,
      text: t('common.name'),
    },
    {
      value: ENTITY_PAYLOADS_VARIABLES.type,
      text: t('common.type'),
    },
    {
      value: ENTITY_PAYLOADS_VARIABLES.infos,
      text: t('common.infos'),
      variables: infosVariables.value,
    },
    {
      value: ENTITY_PAYLOADS_VARIABLES.connector,
      text: t('common.connector'),
    },
    {
      value: ENTITY_PAYLOADS_VARIABLES.component,
      text: t('common.component'),
    },
    {
      value: ENTITY_PAYLOADS_VARIABLES.connectorName,
      text: t('common.connectorName'),
    },
    {
      value: ENTITY_PAYLOADS_VARIABLES.resource,
      text: t('common.resource'),
    },
    {
      value: ENTITY_PAYLOADS_VARIABLES.state,
      text: t('common.state'),
    },
    {
      value: ENTITY_PAYLOADS_VARIABLES.status,
      text: t('common.status'),
    },
    {
      value: ENTITY_PAYLOADS_VARIABLES.stateOutput,
      text: tc('entity.fields.stateOutput'),
    },
    {
      value: ENTITY_PAYLOADS_VARIABLES.impactLevel,
      text: tc('common.impactLevel'),
    },
    {
      value: ENTITY_PAYLOADS_VARIABLES.impactState,
      text: tc('common.impactState'),
    },
    {
      value: ENTITY_PAYLOADS_VARIABLES.category,
      text: tc('common.category'),
    },
  ]);

  return {
    variables,
  };
};
