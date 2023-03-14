<template lang="pug">
  v-layout(column)
    v-layout.mt-3(row, align-end)
      v-flex(xs6)
        div.subheading {{ $t('linkRule.type') }}
        v-radio-group(
          :value="form.type",
          row,
          mandatory,
          @change="updateType"
        )
          v-radio(
            v-for="type in types",
            :value="type.value",
            :label="type.label",
            :key="type.value",
            color="primary"
          )
      v-flex(xs6)
        c-enabled-field(v-field="form.enabled")
    c-name-field.mb-3(v-field="form.name", required)
    c-patterns-field(
      v-field="form.patterns",
      :alarm-attributes="alarmPatternAttributes",
      :entity-attributes="entityPatternAttributes",
      some-required,
      with-alarm,
      with-entity
    )
    c-collapse-panel.my-3(:title="$t('externalData.title')")
      external-data-form(
        v-field="form.external_data",
        :types="externalDataTypes",
        :variables="externalDataPayloadVariables"
      )
</template>

<script>
import {
  ALARM_PATTERN_FIELDS,
  ENTITY_PATTERN_FIELDS,
  QUICK_RANGES,
  EXTERNAL_DATA_TYPES,
  LINK_RULE_TYPES,
  LINK_RULE_TYPES_TO_DEFAULT_SOURCE_CODES,
} from '@/constants';

import { formMixin, formValidationHeaderMixin } from '@/mixins/form';
import { payloadVariablesMixin } from '@/mixins/payload/variables';

import ExternalDataForm from '@/components/forms/external-data/external-data-form.vue';

import LinkRuleLinksForm from './partials/link-rule-links-form.vue';

export default {
  inject: ['$validator'],
  components: {
    ExternalDataForm,
    LinkRuleLinksForm,
  },
  mixins: [
    formMixin,
    formValidationHeaderMixin,
    payloadVariablesMixin,
  ],
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
    isAlarmType() {
      return this.form.type === LINK_RULE_TYPES.alarm;
    },

    types() {
      return Object.values(LINK_RULE_TYPES).map(type => ({
        value: type,
        label: this.$t(`linkRule.types.${type}`),
      }));
    },

    externalDataTypes() {
      return [{
        text: this.$t(`externalData.types.${EXTERNAL_DATA_TYPES.mongo}`),
        value: EXTERNAL_DATA_TYPES.mongo,
      }];
    },

    externalDataPayloadVariables() {
      return this.isAlarmType
        ? this.alarmPayloadVariables
        : this.entityPayloadVariables;
    },

    alarmPatternAttributes() {
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

    entityPatternAttributes() {
      return [
        {
          value: ENTITY_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
      ];
    },
  },
  methods: {
    updateType(type) {
      this.updateModel({
        ...this.form,

        type,
        source_code: LINK_RULE_TYPES_TO_DEFAULT_SOURCE_CODES[type] ?? '',
      });
    },
  },
};
</script>
