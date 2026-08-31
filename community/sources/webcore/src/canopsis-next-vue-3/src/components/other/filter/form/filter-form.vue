<template>
  <c-layout gap="3" column>
    <c-name-field
      v-if="withTitle"
      v-model="title"
      :label="$t('common.title')"
      name="title"
      autofocus
      required
    />

    <c-patterns-field
      v-model:value="patterns"
      :alarm-attributes="alarmAttributes"
      :entity-attributes="entityAttributes"
      :pbehavior-attributes="pbehaviorAttributes"
      :pending="pending"
      :with-alarm="withAlarm"
      :with-entity="withEntity"
      :with-event="withEvent"
      :with-pbehavior="withPbehavior"
      :with-service-weather="withServiceWeather"
      :event-attributes="eventAttributes"
      some-required
    />
  </c-layout>
</template>

<script>
import { pick } from 'lodash';
import { computed } from 'vue';

import { PATTERNS_FIELDS } from '@/constants';

import { useFormContext } from '@/hooks/form';
import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';

const FILTER_FORM_PATTERNS_FIELDS = [
  PATTERNS_FIELDS.alarm,
  PATTERNS_FIELDS.entity,
  PATTERNS_FIELDS.event,
  PATTERNS_FIELDS.pbehavior,
  PATTERNS_FIELDS.serviceWeather,
];

export default {
  props: {
    withAlarm: {
      type: Boolean,
      default: false,
    },
    withEntity: {
      type: Boolean,
      default: false,
    },
    withPbehavior: {
      type: Boolean,
      default: false,
    },
    withEvent: {
      type: Boolean,
      default: false,
    },
    withServiceWeather: {
      type: Boolean,
      default: false,
    },
    withTitle: {
      type: Boolean,
      default: false,
    },
    entityTypes: {
      type: Array,
      required: false,
    },
    entityCountersType: {
      type: Boolean,
      default: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const form = useFormContext();
    const title = form.model('title');
    const { fetchWidgetFilterPatternFields } = usePatternsFields();

    const { pending, alarmAttributes, entityAttributes, eventAttributes, pbehaviorAttributes } =
      usePatternsFieldsFetching(
        () =>
          fetchWidgetFilterPatternFields({
            params: pick(props, ['entityTypes', 'entityCountersType']),
          }),
        props.readonly,
      );

    const patterns = computed({
      get: () => pick(form.values.value, FILTER_FORM_PATTERNS_FIELDS),
      set: (value) => {
        Object.entries(pick(value, FILTER_FORM_PATTERNS_FIELDS)).forEach(([field, fieldValue]) => {
          form.setValue(field, fieldValue);
        });
      },
    });

    return {
      title,
      patterns,
      pending,
      alarmAttributes,
      entityAttributes,
      eventAttributes,
      pbehaviorAttributes,
    };
  },
};
</script>
