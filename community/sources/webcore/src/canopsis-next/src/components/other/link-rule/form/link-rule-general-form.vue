<template>
  <v-layout column>
    <v-layout align-end>
      <v-flex xs6>
        <div class="text-subtitle-1">
          {{ $t('linkRule.type') }}
        </div>
        <v-radio-group
          :value="form.type"
          row
          mandatory
          @change="updateType"
        >
          <v-radio
            v-for="type in types"
            :key="type.value"
            :value="type.value"
            :label="type.label"
            color="primary"
          />
        </v-radio-group>
      </v-flex>
    </v-layout>
    <c-name-field
      v-field="form.name"
      class="mb-3"
      autofocus
      required
    />
    <c-patterns-field
      v-field="form.patterns"
      :alarm-attributes="alarmPatternAttributes"
      :entity-attributes="entityPatternAttributes"
      :with-alarm="isAlarmType"
      some-required
      with-entity
    />
    <c-collapse-panel
      :title="$t('externalData.title')"
      class="my-3"
    >
      <external-data-form
        v-field="form.external_data"
        :types="externalDataTypes"
        :variables="templateVars.external_data"
      />
    </c-collapse-panel>
  </v-layout>
</template>

<script>
import {
  ALARM_PATTERN_FIELDS,
  ENTITY_PATTERN_FIELDS,
  EXTERNAL_DATA_TYPES,
  LINK_RULE_TYPES,
  LINK_RULE_TYPES_TO_DEFAULT_SOURCE_CODES,
} from '@/constants';

import { formMixin, formValidationHeaderMixin } from '@/mixins/form';

import ExternalDataForm from '@/components/forms/external-data/external-data-form.vue';

export default {
  inject: ['$validator'],
  components: { ExternalDataForm },
  mixins: [
    formMixin,
    formValidationHeaderMixin,
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
    templateVars: {
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
        text: this.$t(`externalData.types.${EXTERNAL_DATA_TYPES.table}`),
        value: EXTERNAL_DATA_TYPES.table,
      }];
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
        },
        {
          value: ALARM_PATTERN_FIELDS.creationDate,
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
