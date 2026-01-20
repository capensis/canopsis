<template>
  <c-patterns-field
    v-field="patterns"
    :entity-attributes="entityAttributes"
    :entity-types="entityTypes"
    :pending="pending"
    required
    with-entity
    entity-counters-type
  />
</template>

<script>
import { STATE_SETTING_ENTITY_TYPES } from '@/constants';

import { useValidationHeader } from '@/hooks/validator/validation-header';
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';

export default {
  inject: ['$validator'],
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Object,
      default: () => ({}),
    },
    entityTypes: {
      type: Array,
      default: () => [...STATE_SETTING_ENTITY_TYPES],
    },
  },
  setup() {
    const { fetchStateSettingPatternFields } = usePatternsFields();
    const { hasAnyError } = useValidationHeader();

    const {
      pending,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchStateSettingPatternFields);

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs
       */
      hasAnyError,
      pending,
      entityAttributes,
    };
  },
};
</script>
