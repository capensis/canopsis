<template lang="pug">
  v-layout(column)
    v-text-field(
      v-field="form.name",
      v-validate="'required'",
      :label="$t('common.name')",
      :error-messages="errors.collect('name')",
      name="name"
    )
    c-enabled-field(v-field="form.enabled")
    c-patterns-field(
      v-field="form.patterns",
      :alarm-attributes="alarmAttributes",
      :entity-attributes="entityAttributes",
      some-required,
      with-alarm,
      with-entity
    )
    c-collapse-panel.my-3(:title="$t('externalData.title')")
      external-data-form(
        v-field="form.external_data",
        :types="externalDataTypes"
      )
</template>

<script>
import {
  ALARM_PATTERN_FIELDS,
  ENTITY_PATTERN_FIELDS,
  QUICK_RANGES,
  EXTERNAL_DATA_TYPES,
} from '@/constants';

import { formValidationHeaderMixin } from '@/mixins/form/validation-header';

import ExternalDataForm from '@/components/forms/external-data/external-data-form.vue';

import LinkRuleLinksForm from './partials/link-rule-links-form.vue';

export default {
  inject: ['$validator'],
  components: {
    ExternalDataForm,
    LinkRuleLinksForm,
  },
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    externalDataTypes() {
      return [{
        text: this.$t(`externalData.types.${EXTERNAL_DATA_TYPES.mongo}`),
        value: EXTERNAL_DATA_TYPES.mongo,
      }];
    },

    alarmAttributes() {
      return [
        {
          value: ALARM_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.lastUpdateDate,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.resolved,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.ackAt,
          options: {
            intervalRanges: [QUICK_RANGES.custom],
          },
        },
        {
          value: ALARM_PATTERN_FIELDS.creationDate,
          options: {
            intervalRanges: [QUICK_RANGES.custom],
          },
        },
      ];
    },

    entityAttributes() {
      return [
        {
          value: ENTITY_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
      ];
    },
  },
};
</script>
