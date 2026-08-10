<template>
  <c-layout gap="2" column>
    <c-text-field
      v-model="titleField"
      :label="$t('common.title')"
      name="title"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="$t('modals.createDynamicInfoTemplate.fields.names')" indented>
        <dynamic-info-template-names-field>
          <template #no-data>
            <v-alert type="info">
              {{ $t('modals.createDynamicInfoTemplate.emptyNames') }}
            </v-alert>
          </template>
        </dynamic-info-template-names-field>
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { useFormModel } from '@/hooks/form';

import DynamicInfoTemplateNamesField from './fields/dynamic-info-template-names-field.vue';

export default {
  components: { DynamicInfoTemplateNamesField },

  setup() {
    const titleField = useFormModel('title');

    return {
      titleField,
    };
  },
};
</script>
