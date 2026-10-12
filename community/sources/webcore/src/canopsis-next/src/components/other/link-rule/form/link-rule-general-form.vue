<template>
  <v-layout class="gap-3" column>
    <c-name-field
      v-field="form.name"
      :max-length="255"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="$t('linkRule.type')">
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
      </c-form-block-row>

      <c-form-block-row :label="$t('externalData.title')">
        <external-data-form
          v-field="form.external_data"
          :types="externalDataTypes"
          :variables="templateVars.external_data"
        />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { EXTERNAL_DATA_TYPES, LINK_RULE_TYPES, LINK_RULE_TYPES_TO_DEFAULT_SOURCE_CODES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useValidationHeader } from '@/hooks/validator/validation-header';

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
    const { hasAnyError } = useValidationHeader();

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
      isAlarmType,
      types,
      externalDataTypes,
      updateType,
    };
  },
};
</script>
