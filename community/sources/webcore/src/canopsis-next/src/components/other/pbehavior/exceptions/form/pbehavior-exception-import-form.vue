<template>
  <c-layout gap="3" column>
    <c-name-field v-model="name" autofocus required />

    <c-form-block>
      <c-form-block-row :label="t('pbehavior.pbehaviorType')">
        <c-pbehavior-type-field v-model="type" name="type" required />
      </c-form-block-row>

      <c-form-block-row :label="t('common.file')" align-center>
        <file-selector
          v-model="file"
          name="file"
          required
          with-files-list
          hide-details
          @change="changeFile"
        >
          <template #activator="{ on, disabled, errors }">
            <c-btn-with-error
              :error="errors[0]"
              :disabled="disabled"
              variant="outlined"
              v-on="on"
            >
              {{ t('common.chooseFile') }}
            </c-btn-with-error>
          </template>
        </file-selector>
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import FileSelector from '@/components/forms/fields/file-selector.vue';

export default {
  components: { FileSelector },

  setup() {
    const { t } = useI18n();
    const { name, type, file } = useFormModels(['name', 'type', 'file']);

    const changeFile = selectedFile => file.value = selectedFile;

    return {
      t,
      name,
      type,
      file,
      changeFile,
    };
  },
};
</script>
