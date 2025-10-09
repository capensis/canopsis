<template>
  <v-layout column>
    <span class="text-body-2 my-2">{{ $tc('common.column', 2) }}</span>
    <v-flex xs12>
      <v-alert
        :value="!form.columns.length"
        color="info"
      >
        {{ $t('widgetTemplate.errors.columnsRequired') }}
      </v-alert>
    </v-flex>
    <c-columns-field
      v-field="form.columns"
      :type="entityType"
      with-color-indicator
      with-template
      with-html
      @input="validate"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { ENTITIES_TYPES, WIDGET_TEMPLATES_TYPES } from '@/constants';

import { useValidationAttachRequiredForField } from '@/hooks/validator/validation-attach-required';

export default {
  inject: ['$validator'],
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

  setup(props) {
    const name = 'columns';

    const {
      asyncValidateRequiredRule: validate,
    } = useValidationAttachRequiredForField(name, () => !!props.form.columns?.length);

    const entityType = computed(() => (props.form.type === WIDGET_TEMPLATES_TYPES.alarmColumns
      ? ENTITIES_TYPES.alarm
      : ENTITIES_TYPES.entity));

    return {
      validate,
      entityType,
    };
  },
};
</script>
