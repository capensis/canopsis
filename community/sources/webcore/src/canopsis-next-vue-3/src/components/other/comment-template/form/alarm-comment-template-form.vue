<template>
  <v-card class="c-alternative-bg-panel" rounded="xl">
    <v-card-text>
      <c-layout gap="2" column>
        <c-select
          v-if="templates.length"
          :model-value="form.template"
          :items="templates"
          :label="t('common.template', 1)"
          name="template"
          item-title="name"
          item-value="_id"
          return-object
          clearable
          @update:model-value="updateField('template', $event)"
        />
        <c-name-field
          v-for="(field, index) in fields"
          :key="field.name"
          :model-value="form[field.name]"
          :required="field.required"
          :label="field.label || field.name"
          :name="field.name"
          :max-length="FORM_MAX_STRING_LENGTH"
          :autofocus="index === 0"
          @update:model-value="updateField(field.name, $event)"
        />
      </c-layout>
    </v-card-text>
  </v-card>
</template>

<script>
import { computed } from 'vue';

import { FORM_MAX_STRING_LENGTH } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

export default {
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

  emits: ['update:form'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    const fields = computed(() => (
      props.form.template?.fields ?? [{ name: 'comment', label: t('common.note'), required: true }]
    ));

    return {
      t,
      FORM_MAX_STRING_LENGTH,
      fields,
      updateField,
    };
  },
};
</script>
