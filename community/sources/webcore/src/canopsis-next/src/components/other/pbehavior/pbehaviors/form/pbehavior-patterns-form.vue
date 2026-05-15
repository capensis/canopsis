<template>
  <c-patterns-field
    v-field="form"
    :entity-attributes="entityAttributes"
    :readonly="readonly"
    :pending="pending"
    :counter-method="counterMethod"
    entity-counters-type
    with-entity
    required
  />
</template>

<script>
import { computed } from 'vue';

import { useStoreModuleHooks } from '@/hooks/store';
import { usePatternsFields, usePatternsFieldsFetching } from '@/hooks/store/modules/patterns-fields';

const usePbehaviorPatternsStoreModule = () => useStoreModuleHooks('pbehaviorPatterns');

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
  setup(props) {
    const { fetchPbehaviorPatternFields } = usePatternsFields();

    const {
      pending,
      entityAttributes,
    } = usePatternsFieldsFetching(fetchPbehaviorPatternFields, props.readonly);

    const { useActions: usePbehaviorPatternActions } = usePbehaviorPatternsStoreModule();
    const { checkPatternsPbehaviorsCount } = usePbehaviorPatternActions({
      checkPatternsPbehaviorsCount: 'checkPatternsPbehaviorsCount',
    });

    const checkFilter = async ({ data } = {}) => {
      const counter = await checkPatternsPbehaviorsCount({
        data: { ...data, _id: props.pbehaviorId },
      });

      return {
        entity_pattern: counter,
        all: counter,
      };
    };

    const counterMethod = computed(() => (
      props.pbehaviorCounterType
        ? checkFilter
        : undefined
    ));

    return {
      pending,
      entityAttributes,
      counterMethod,
    };
  },
};
</script>
