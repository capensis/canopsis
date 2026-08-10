<template>
  <c-patterns-field
    :value="patterns"
    :entity-attributes="entityAttributes"
    :entity-types="entityTypes"
    :pending="pending"
    required
    with-entity
    entity-counters-type
    @update:value="patterns = $event"
  />
</template>

<script>
import { STATE_SETTING_ENTITY_TYPES } from '@/constants';
import { useFormModel } from '@/hooks/form';
import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';

export default {
  props: {
    name: {
      type: String,
      default: 'entity_pattern',
    },
    entityTypes: {
      type: Array,
      default: () => [...STATE_SETTING_ENTITY_TYPES],
    },
  },

  setup(props) {
    const patterns = useFormModel(props.name);
    const { fetchStateSettingPatternFields } = usePatternsFields();

    const { pending, entityAttributes } = usePatternsFieldsFetching(fetchStateSettingPatternFields);

    return {
      patterns,
      pending,
      entityAttributes,
    };
  },
};
</script>
