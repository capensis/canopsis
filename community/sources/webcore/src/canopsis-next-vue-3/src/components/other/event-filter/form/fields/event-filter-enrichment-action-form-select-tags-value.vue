<template>
  <c-select
    :model-value="value"
    :label="$t('common.value')"
    :items="items"
    :name="name"
    class="ml-2"
    required
    @update:model-value="updateValue"
  />
</template>

<script>
import { isEqual } from 'lodash';
import { computed, nextTick, onBeforeUnmount, onMounted, watch } from 'vue';

import { EVENT_FILTER_SET_TAGS_REGEX } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useValidationAttachRule } from '@/hooks/validator/validation-attach-rule';

export default {
  emits: ['update:value'],
  props: {
    value: {
      type: String,
      required: false,
    },
    items: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'tags_value',
    },
  },
  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);

    const selectedItem = computed(() => props.items.find(({ value }) => value === props.value));
    const regexRuleFieldName = computed(() => `${props.name}.regex`);

    const {
      attachRule: attachRegexRule,
      detachRule: detachRegexRule,
      validateRule: validateRegexRule,
    } = useValidationAttachRule(
      regexRuleFieldName,
      {
        regex: EVENT_FILTER_SET_TAGS_REGEX,
      },
      () => selectedItem.value?.valueForValidation ?? '',
    );

    const validateRegex = () => nextTick(validateRegexRule);

    const updateValue = (value) => {
      updateModel(value);
      validateRegex();
    };

    watch(
      () => props.items,
      (items, oldItems) => {
        if (props.value && !isEqual(items, oldItems) && items.every(item => item.value !== props.value)) {
          updateModel('');
        }

        validateRegex();
      },
      { immediate: true },
    );

    onMounted(attachRegexRule);
    onBeforeUnmount(detachRegexRule);

    return {
      updateValue,
    };
  },
};
</script>
