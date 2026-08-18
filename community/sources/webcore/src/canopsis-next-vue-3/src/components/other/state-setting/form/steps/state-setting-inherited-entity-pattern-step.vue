<template>
  <c-patterns-field
    :value="patterns"
    :entity-attributes="entityAttributes"
    :entity-types="entityTypes"
    :entity-title="t('stateSetting.dependenciesEntityPattern')"
    :disabled="disabled"
    :pending="pending"
    entity-name="inherited_entity_pattern"
    required
    with-entity
    entity-counters-type
    @update:value="patterns = $event"
  />
</template>

<script>
import { computed } from 'vue';

import { STATE_SETTING_ENTITY_TYPES } from '@/constants';

import { useFormModel } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';
import { useValidationHeader } from '@/hooks/validator/validation-header';

export default {
  props: {
    name: {
      type: String,
      default: 'inherited_entity_pattern',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const patterns = useFormModel(props.name);
    const { fetchStateSettingPatternFields } = usePatternsFields();
    const { hasAnyError } = useValidationHeader();

    const { pending, entityAttributes } = usePatternsFieldsFetching(
      fetchStateSettingPatternFields,
      props.disabled,
    );

    const entityTypes = computed(() => [...STATE_SETTING_ENTITY_TYPES]);

    return {
      t,
      /**
       * It's using in the parent component to display the validation header color for tabs
       */
      hasAnyError,
      patterns,
      pending,
      entityAttributes,
      entityTypes,
    };
  },
};
</script>
