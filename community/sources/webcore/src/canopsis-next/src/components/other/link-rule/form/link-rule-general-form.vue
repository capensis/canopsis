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
      <v-flex xs6>
        <c-enabled-field v-field="form.enabled" />
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
      :alarm-attributes="alarmAttributes"
      :entity-attributes="entityAttributes"
      :pending="pending"
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
import { computed } from 'vue';

import { EXTERNAL_DATA_TYPES, LINK_RULE_TYPES, LINK_RULE_TYPES_TO_DEFAULT_SOURCE_CODES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useValidationHeader } from '@/hooks/validator/validation-header';
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';

import ExternalDataForm from '@/components/forms/external-data/external-data-form.vue';

export default {
  inject: ['$validator'],
  components: { ExternalDataForm },
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
  setup(props, { emit }) {
    const { t } = useI18n();
    const { fetchLinkRulePatternFields } = usePatternsFields();
    const { hasAnyError } = useValidationHeader();

    const {
      pending,
      alarmAttributes,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchLinkRulePatternFields);

    const isAlarmType = computed(() => props.form.type === LINK_RULE_TYPES.alarm);

    const types = computed(() => Object.values(LINK_RULE_TYPES).map(type => ({
      value: type,
      label: t(`linkRule.types.${type}`),
    })));

    const externalDataTypes = computed(() => [{
      text: t(`externalData.types.${EXTERNAL_DATA_TYPES.table}`),
      value: EXTERNAL_DATA_TYPES.table,
    }]);

    const updateType = type => emit('input', {
      ...props.form,
      type,
      source_code: LINK_RULE_TYPES_TO_DEFAULT_SOURCE_CODES[type] ?? '',
    });

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs
       */
      hasAnyError,
      pending,
      alarmAttributes,
      entityAttributes,
      isAlarmType,
      types,
      externalDataTypes,
      updateType,
    };
  },
};
</script>
