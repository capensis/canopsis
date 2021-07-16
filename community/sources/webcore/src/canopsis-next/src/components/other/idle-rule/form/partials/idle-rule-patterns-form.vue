<template lang="pug">
  div
    patterns-list(
      v-if="isEntityType",
      v-field="form.entity_patterns",
      v-validate="'required'",
      name="entity_patterns"
    )
    c-patterns-field(
      v-else,
      v-field="form",
      some-required,
      alarm,
      entity
    )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { PATTERNS_TYPES } from '@/constants';

import { formValidationHeaderMixin } from '@/mixins/form';

import PatternsList from '@/components/common/patterns-list/patterns-list.vue';

const { mapActions } = createNamespacedHelpers('idleRules');

export default {
  provide() {
    return {
      $checkEntitiesCountByType: this.checkEntitiesCountByType,
    };
  },
  inject: ['$validator'],
  components: { PatternsList },
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
    isEntityType: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    ...mapActions({
      fetchIdleRuleEntitiesCountWithoutStore: 'fetchEntitiesCountWithoutStore',
    }),

    async checkEntitiesCountByType(type, patterns) {
      const responseKey = PATTERNS_TYPES.alarm ? 'total_count_alarms' : 'total_count_entities';

      const result = await this.fetchIdleRuleEntitiesCountWithoutStore({
        data: { [`${type}_patterns`]: patterns },
      });

      return {
        over_limit: result.over_limit,
        total_count: result[responseKey],
      };
    },
  },
};
</script>
