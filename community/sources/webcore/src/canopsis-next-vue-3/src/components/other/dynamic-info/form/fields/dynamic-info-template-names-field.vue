<template>
  <c-layout gap="2" column>
    <slot v-if="!namesField.length" name="no-data" />

    <c-layout
      v-for="(name, index) in namesField"
      :key="name.key"
      align="center"
      gap="2"
    >
      <c-name-field v-model="name.value" :name="`names[${index}].value`" required />

      <c-action-btn type="delete" @click="remove(index)" />
    </c-layout>

    <c-layout>
      <v-btn color="primary" variant="outlined" @click="addName">
        {{ $t('modals.createDynamicInfoTemplate.buttons.addName') }}
      </v-btn>
    </c-layout>

    <v-alert v-if="hasNamesError" type="error">
      {{ $t('modals.createDynamicInfoTemplate.errors.noNames') }}
    </v-alert>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { generateTemplateFormName } from '@/helpers/entities/dynamic-info/template/form';

import { useFormContext } from '@/hooks/form';

export default {
  setup() {
    const form = useFormContext();
    const namesField = form.array('names');
    const hasNamesError = computed(() => form.submitCount.value > 0 && !namesField.value.length);

    const addName = () => namesField.value.push(generateTemplateFormName());
    const remove = index => namesField.value.splice(index, 1);

    return {
      namesField,
      hasNamesError,
      addName,
      remove,
    };
  },
};
</script>
