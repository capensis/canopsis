<template>
  <c-layout gap="3" column>
    <c-name-field v-model="name" autofocus required />

    <c-form-block>
      <c-form-block-row :label="t('pbehavior.pbehaviorType')">
        <c-pbehavior-type-field v-model="type" name="type" required />
      </c-form-block-row>

      <c-form-block-row :label="t('common.file')" align-center>
        <file-selector v-model="file" name="file" required>
          <template #activator="{ on, disabled }">
            <v-btn
              :disabled="disabled"
              color="primary"
              variant="outlined"
              v-on="on"
            >
              {{ t('common.chooseFile') }}
            </v-btn>
          </template>
        </file-selector>
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
    const { name, type, file } = useFormModels(['name', 'type', 'file']);

    return {
      t,
      name,
      type,
      file,
    };
  },
};
</script>
