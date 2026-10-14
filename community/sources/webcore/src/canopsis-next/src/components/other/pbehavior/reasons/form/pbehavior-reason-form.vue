<template>
  <c-layout gap="3" column>
    <c-enabled-field v-model="visible" :label="t('pbehavior.visible')" with-background />

    <c-name-field
      v-model="name"
      :label="t('modals.createPbehaviorReason.fields.name')"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="t('modals.createPbehaviorReason.fields.description')" indented>
        <c-description-field v-model="description" name="description" required />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  setup() {
    const { t } = useI18n();
    const { name, description, visible } = useFormModels(['name', 'description', 'visible']);

    return {
      t,
      name,
      description,
      visible,
    };
  },
};
</script>
