<template>
  <v-layout column>
    <span class="text-body-2 my-2">{{ $tc('common.column', 2) }}</span>
    <v-flex xs12>
      <v-alert
        :value="required && !value.length"
        color="info"
      >
        {{ $t('widgetTemplate.errors.sortColumnsRequired') }}
      </v-alert>
    </v-flex>
    <c-sort-columns-field
      v-field="value"
      :items="items"
      @input="validate"
    />
  </v-layout>
</template>

<script>
import { watch, nextTick, onBeforeUnmount } from 'vue';

import { useValidationAttachRequired } from '@/hooks/validator/validation-attach-required';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
    required: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'sort_columns',
    },
    items: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const { attachRequiredRule, detachRequiredRule, validateRequiredRule } = useValidationAttachRequired(props.name);

    const getter = () => !!props.value.length;
    const validate = () => props.required && nextTick(validateRequiredRule);

    watch(() => props.required, (required) => {
      if (required) {
        attachRequiredRule(getter);
        validateRequiredRule();

        return;
      }

      detachRequiredRule();
    }, { immediate: true });

    onBeforeUnmount(detachRequiredRule);

    return {
      validate,
    };
  },
};
</script>
