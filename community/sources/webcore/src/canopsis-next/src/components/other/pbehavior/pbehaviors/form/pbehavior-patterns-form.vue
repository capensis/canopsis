<template>
  <c-patterns-field
    :value="form"
    :entity-attributes="entityAttributes"
    :readonly="readonly"
    :pending="pending"
    :counter-method="counterMethod"
    entity-counters-type
    with-entity
    required
    @update:value="$emit('update:form', $event)"
  />
</template>

<script>
import { computed } from 'vue';

import {
  usePatternsFields,
  usePatternsFieldsFetching,
} from '@/hooks/store/modules/patterns-fields';
import { usePbehaviorPatterns } from '@/hooks/store/modules/pbehavior-patterns';

export default {
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    pbehaviorId: {
      type: String,
      required: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    pbehaviorCounterType: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:form'],

  setup(props) {
    const { fetchPbehaviorPatternFields } = usePatternsFields();

    const { pending, entityAttributes } = usePatternsFieldsFetching(
      fetchPbehaviorPatternFields,
      props.readonly,
    );

    const { checkPatternsPbehaviorsCount } = usePbehaviorPatterns();

    const checkFilter = async ({ data } = {}) => {
      const counter = await checkPatternsPbehaviorsCount({
        data: { ...data, _id: props.pbehaviorId },
      });

      return {
        entity_pattern: counter,
        combined: counter,
      };
    };

    const counterMethod = computed(() => (props.pbehaviorCounterType ? checkFilter : undefined));

    return {
      pending,
      entityAttributes,
      counterMethod,
    };
  },
};
</script>
