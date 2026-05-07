<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      :max-length="255"
      autofocus
      required
    />
    <c-form-block>
      <c-form-block-row :label="$t('common.priority')">
        <c-priority-field v-field="form.priority" :disabled="defaultRule" />
      </c-form-block-row>
      <c-form-block-row :label="$t('common.duration')">
        <c-duration-field v-field="form.duration" required />
      </c-form-block-row>
      <c-form-block-row
        v-if="flapping"
        :label="$t('common.frequencyLimit')"
      >
        <c-number-field
          v-field="form.freq_limit"
          :label="$t('common.frequencyLimit')"
          :min="1"
          name="freq_limit"
        />
      </c-form-block-row>
      <c-form-block-row :label="$t('common.description')">
        <c-description-field v-field="form.description" required />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { useValidationHeader } from '@/hooks/validator/validation-header';

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
    flapping: {
      type: Boolean,
      default: false,
    },
    defaultRule: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { hasAnyError } = useValidationHeader();

    return {
      hasAnyError,
    };
  },
};
</script>
