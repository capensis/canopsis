<template>
  <c-patterns-field
    v-field="form"
    :entity-attributes="entityAttributes"
    :readonly="readonly"
    :counter-method="counterMethod"
    entity-counters-type
    with-entity
    required
  />
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { ENTITY_PATTERN_FIELDS } from '@/constants';

import { formValidationHeaderMixin } from '@/mixins/form';

const { mapActions: mapPbehaviorPatternActions } = createNamespacedHelpers('pbehaviorPatterns');

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
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
  computed: {
    entityAttributes() {
      return [
        {
          value: ENTITY_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
      ];
    },

    counterMethod() {
      return this.pbehaviorCounterType
        ? this.checkFilter
        : undefined;
    },
  },
  methods: {
    ...mapPbehaviorPatternActions({
      checkPatternsPbehaviorsCount: 'checkPatternsPbehaviorsCount',
    }),

    async checkFilter({ data } = {}) {
      const counter = await this.checkPatternsPbehaviorsCount({
        data: { ...data, _id: this.pbehaviorId },
      });

      return {
        entity_pattern: counter,
        all: counter,
      };
    },
  },
};
</script>
