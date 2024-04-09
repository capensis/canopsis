<template>
  <v-layout column>
    <v-layout>
      <v-flex xs6>
        <c-number-field
          v-field="config.threshold_count"
          :min="0"
          :label="$t('metaAlarmRule.thresholdCount')"
          name="thresholdCount"
        >
          <template #append-outer>
            <c-help-icon
              :text="$t('metaAlarmRule.thresholdCountHelpText')"
              max-width="300"
              icon="help"
              top
            />
          </template>
        </c-number-field>
      </v-flex>
    </v-layout>
    <meta-alarm-rule-time-based-form v-field="config" with-child-inactive-delay />
    <c-payload-text-field
      v-field="config.corel_id"
      :label="$t('metaAlarmRule.corelId')"
      name="corelId"
      required
    >
      <template #append-outer="">
        <c-help-icon
          :text="$t('metaAlarmRule.corelIdHelpText')"
          max-width="300"
          icon="help"
          left
        />
      </template>
    </c-payload-text-field>
    <c-payload-text-field
      v-field="config.corel_status"
      :label="$t('metaAlarmRule.corelStatus')"
      name="corelStatus"
      required
    >
      <template #append-outer="">
        <c-help-icon
          :text="$t('metaAlarmRule.corelStatusHelpText')"
          max-width="300"
          icon="help"
          left
        />
      </template>
    </c-payload-text-field>
    <c-payload-text-field
      v-field="config.corel_parent"
      :label="$t('metaAlarmRule.corelParent')"
      name="corelParent"
      required
    >
      <template #append-outer="">
        <c-help-icon
          :text="$t('metaAlarmRule.corelParentHelpText')"
          max-width="300"
          icon="help"
          left
        />
      </template>
    </c-payload-text-field>
    <c-payload-text-field
      v-field="config.corel_child"
      :label="$t('metaAlarmRule.corelChild')"
      :variables="variables"
      name="corelChild"
      required
    >
      <template #append-outer="">
        <c-help-icon
          :text="$t('metaAlarmRule.corelChildHelpText')"
          max-width="300"
          icon="help"
          left
        />
      </template>
    </c-payload-text-field>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { ALARM_PAYLOADS_VARIABLES, ENTITY_PAYLOADS_VARIABLES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import MetaAlarmRuleTimeBasedForm from './meta-alarm-rule-time-based-form.vue';

export default {
  components: { MetaAlarmRuleTimeBasedForm },
  model: {
    prop: 'config',
    event: 'input',
  },
  props: {
    config: {
      type: Object,
      default: () => ({}),
    },
    entityInfos: {
      type: Array,
      default: () => [],
    },
  },
  setup() {
    const { t, tc } = useI18n();

    const entityPayloadVariables = computed(() => []);
    const alarmPayloadVariables = computed(() => [
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

    const variables = computed(() => [
      {
        value: ENTITY_PAYLOADS_VARIABLES.entity,
        text: tc('common.entity'),
        variables: entityPayloadVariables.value,
      },
      {
        value: ALARM_PAYLOADS_VARIABLES.alarm,
        text: tc('common.alarm'),
        variables: alarmPayloadVariables.value,
      },
    ]);

    return {
      variables,
    };
  },
};
</script>
