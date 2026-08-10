<template>
  <c-patterns-field
    v-model:value="patternsField"
    :some-required="!isChangeEntityType"
    :required="isChangeEntityType"
    :with-entity="!isChangeEntityType"
    :event-attributes="eventAttributes"
    :pending="attributesPending"
    expanded-event
    with-event
    entity-counters-type
  />
</template>

<script>
import { computed } from 'vue';

import { isChangeEntityEventFilterRuleType } from '@/helpers/entities/event-filter/rule/entity';
import { useFormModel } from '@/hooks/form';

export default {
  emits: ['update:form'],
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    eventAttributes: {
      type: Array,
      required: false,
    },
    attributesPending: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const patternsField = useFormModel('patterns');

    const isChangeEntityType = computed(() => isChangeEntityEventFilterRuleType(props.form.type));

    return {
      patternsField,

      isChangeEntityType,
    };
  },
};
</script>
