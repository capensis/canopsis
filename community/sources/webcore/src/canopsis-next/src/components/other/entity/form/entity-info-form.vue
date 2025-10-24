<template>
  <v-layout column>
    <v-text-field
      v-field="form.name"
      v-validate="nameRules"
      :label="$t('common.name')"
      :error-messages="errors.collect('name')"
      name="name"
      autofocus
    />
    <v-text-field
      v-field="form.description"
      v-validate="descriptionRules"
      :label="$t('common.description')"
      :error-messages="errors.collect('description')"
      name="description"
    />
    <c-mixed-field
      v-field="form.value"
      :label="$t('common.value')"
      :types="fieldTypes"
      required
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { PATTERN_FIELD_TYPES } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
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
  setup(props) {
    const infosNames = computed(() => props.infos.map(({ name }) => name));

    const descriptionRules = computed(() => ({
      required: true,
    }));

    const nameRules = computed(() => ({
      required: true,
      unique: {
        values: infosNames.value,
        initialValue: props.entityInfo?.name,
      },
    }));

    const fieldTypes = computed(() => [
      { value: PATTERN_FIELD_TYPES.string },
      { value: PATTERN_FIELD_TYPES.number },
      { value: PATTERN_FIELD_TYPES.boolean },
      { value: PATTERN_FIELD_TYPES.stringArray },
    ]);

    return {
      descriptionRules,
      nameRules,
      fieldTypes,
    };
  },
};
</script>
