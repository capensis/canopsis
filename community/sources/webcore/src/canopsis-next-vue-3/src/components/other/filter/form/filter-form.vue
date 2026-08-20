<template>
  <c-layout gap="3" column>
    <c-name-field
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
      :with-pbehavior="withPbehavior"
      :with-service-weather="withServiceWeather"
      some-required
    />
  </c-layout>
</template>

<script>
import { pick } from 'lodash';
import { computed } from 'vue';

import { PATTERNS_FIELDS } from '@/constants';

import { useFormContext, useFormModel } from '@/hooks/form';
import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';

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
    withServiceWeather: {
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
    const title = useFormModel('title');
    const { fetchWidgetFilterPatternFields } = usePatternsFields();

    const { pending, alarmAttributes, entityAttributes, pbehaviorAttributes } =
      usePatternsFieldsFetching(
        () =>
          fetchWidgetFilterPatternFields({
            params: pick(props, ['entityTypes', 'entityCountersType']),
          }),
        props.readonly,
      );

    const patterns = computed({
      get: () =>
        pick(form.values.value, [
          PATTERNS_FIELDS.alarm,
          PATTERNS_FIELDS.entity,
          PATTERNS_FIELDS.pbehavior,
          PATTERNS_FIELDS.serviceWeather,
        ]),
      set: (value) => {
        form.values.value = {
          ...form.values.value,
          ...value,
        };
      },
    });

    return {
      title,
      patterns,
      pending,
      alarmAttributes,
      entityAttributes,
      pbehaviorAttributes,
    };
  },
};
</script>
