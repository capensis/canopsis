<template>
  <c-layout gap="3" column>
    <c-enabled-field v-model="visible" :label="t('pbehavior.visible')" with-background />

    <c-name-field
      v-model="name"
      :label="t('modals.createPbehaviorException.fields.name')"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="t('modals.createPbehaviorException.fields.description')">
        <c-description-field v-model="description" name="description" required />
      </c-form-block-row>

      <c-form-block-row
        :label="t('pbehavior.exdates.title')"
        :align-center="!exdates.length"
        indented
      >
        <pbehavior-exceptions-field
          :exdates="exdates"
          with-exdate-type
          required
          @update:exdates="exdates = $event"
        />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { useFormArray, useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import PbehaviorExceptionsField from '@/components/other/pbehavior/exceptions/fields/pbehavior-exceptions-field.vue';

export default {
  components: { PbehaviorExceptionsField },

  setup() {
    const { t } = useI18n();
    const { name, description, visible } = useFormModels(['name', 'description', 'visible']);
    const exdates = useFormArray('exdates');

    return {
      t,
      name,
      description,
      visible,
      exdates,
    };
  },
};
</script>
