<template>
  <c-layout column>
    <c-name-field
      :model-value="form.name"
      :unique="uniqueNames"
      name="name"
      autofocus
      required
      @update:model-value="updateField('name', $event)"
    />
    <c-description-field
      :model-value="form.description"
      required
      @update:model-value="updateField('description', $event)"
    />
    <c-mixed-field
      :value="form.value"
      :label="$t('common.value')"
      :types="fieldTypes"
      required
      @update:value="updateField('value', $event)"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { PATTERN_FIELD_TYPES } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';

export default {
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    entityInfo: {
      type: Object,
      default: () => ({}),
    },
    infos: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    const uniqueNames = computed(() => props.infos
        .map(({ name }) => name)
        .filter(name => name !== props.entityInfo?.name),
    );

    const fieldTypes = computed(() => [
      { value: PATTERN_FIELD_TYPES.string },
      { value: PATTERN_FIELD_TYPES.number },
      { value: PATTERN_FIELD_TYPES.boolean },
      { value: PATTERN_FIELD_TYPES.stringArray },
    ]);

    return {
      uniqueNames,
      fieldTypes,
      updateField,
    };
  },
};
</script>
