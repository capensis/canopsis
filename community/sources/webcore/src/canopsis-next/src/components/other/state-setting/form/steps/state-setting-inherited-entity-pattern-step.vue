<template>
  <c-patterns-field
    v-field="patterns"
    :entity-attributes="entityAttributes"
    :entity-types="entityTypes"
    :entity-title="$t('stateSetting.dependenciesEntityPattern')"
    :disabled="disabled"
    :pending="pending"
    entity-name="inherited_entity_pattern"
    required
    with-entity
    entity-counters-type
  />
</template>

<script>
import { computed } from 'vue';

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
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { fetchStateSettingPatternFields } = usePatternsFields();
    const { hasAnyError } = useValidationHeader();

    const {
      pending,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchStateSettingPatternFields, props.disabled);

    const entityTypes = computed(() => [...STATE_SETTING_ENTITY_TYPES]);

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs
       */
      hasAnyError,
      pending,
      entityAttributes,
      entityTypes,
    };
  },
};
</script>
