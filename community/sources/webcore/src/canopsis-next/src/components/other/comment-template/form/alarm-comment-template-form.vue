<template>
  <v-card class="c-alternative-bg-panel" rounded="xl">
    <v-card-text>
      <v-layout class="gap-2" column>
        <c-select-field
          v-field="form.template"
          :items="templates"
          :label="$tc('common.template', 1)"
          name="template"
          return-object
          clearable
          @input="clearErrors"
        />
        <c-name-field
          v-for="field in fields"
          v-field="form[field.name]"
          :key="field.name"
          :required="field.required"
          :label="field.label || field.name"
          :name="field.name"
          :max-length="255"
        />
      </v-layout>
    </v-card-text>
  </v-card>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useValidator } from '@/hooks/validator/validator';

export default {
  inject: ['$validator'],
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    templates: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const { t } = useI18n();
    const validator = useValidator();

    const fields = computed(() => (
      props.form.template?.fields ?? [{ name: 'comment', label: t('common.note'), required: true }]
    ));

    const clearErrors = () => validator.errors.clear();

    return {
      fields,
      clearErrors,
    };
  },
};
</script>
