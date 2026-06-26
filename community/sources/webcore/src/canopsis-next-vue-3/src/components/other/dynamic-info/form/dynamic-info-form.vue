<template>
  <c-layout
    gap="2"
    column
  >
    <c-enabled-field
      v-model="enabledField"
      :label="$t('common.enabled')"
      with-background
    />

    <c-name-field
      v-model="nameField"
      name="name"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="$t('common.description')">
        <c-description-field
          v-model="descriptionField"
          name="description"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('modals.createDynamicInfo.infosSection.title')">
        <dynamic-info-infos-form />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { useFormModels } from '@/hooks/form';

import DynamicInfoInfosForm from './dynamic-info-infos-form.vue';

export default {
  components: {
    DynamicInfoInfosForm,
  },
  setup() {
    const {
      enabledField,
      nameField,
      descriptionField,
    } = useFormModels({
      enabledField: 'enabled',
      nameField: 'name',
      descriptionField: 'description',
    });

    return {
      enabledField,
      nameField,
      descriptionField,
    };
  },
};
</script>
