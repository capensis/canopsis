<template>
  <c-patterns-field
    v-field="form.patterns"
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

import { useValidationHeader } from '@/hooks/validator/validation-header';

export default {
  model: {
    prop: 'form',
    event: 'input',
  },
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
    const { hasAnyError } = useValidationHeader();

    const isChangeEntityType = computed(() => isChangeEntityEventFilterRuleType(props.form.type));

    return {
      /**
       * It's using in the parent component to display the validation header color for tabs
       */
      hasAnyError,
      isChangeEntityType,
    };
  },
};
</script>
