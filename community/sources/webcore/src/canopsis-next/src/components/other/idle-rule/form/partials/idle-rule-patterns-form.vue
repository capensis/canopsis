<template lang="pug">
  div
    patterns-list(
      v-if="isEntityType",
      v-field="form.entity_patterns",
      v-validate="'required'",
      :type="$constants.PATTERNS_TYPES.entity",
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

import { formValidationHeaderMixin, validationErrorsMixinCreator } from '@/mixins/form';

import PatternsList from '@/components/common/patterns-list/patterns-list.vue';

const { mapActions } = createNamespacedHelpers('idleRules');

export default {
  provide() {
    return {
      $checkEntitiesCountForPatternsByType: this.checkEntitiesCountForPatternsByType,
    };
  },
  inject: ['$validator'],
  components: { PatternsList },
  mixins: [
    formValidationHeaderMixin,
    validationErrorsMixinCreator(),
  ],
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

    setFormErrors(err) {
      const existFieldErrors = this.getExistsFieldsErrors(err);

      if (existFieldErrors.length) {
        this.addExistsFieldsErrors(existFieldErrors);

        return {
          over_limit: false,
          total_count: 0,
        };
      }

      throw err;
    },

    async checkEntitiesCountForPatternsByType(type, patterns) {
      const requestKey = `${type}_patterns`;
      const responseKey = type === PATTERNS_TYPES.alarm ? 'total_count_alarms' : 'total_count_entities';

      try {
        const response = await this.fetchIdleRuleEntitiesCountWithoutStore({
          data: { [requestKey]: patterns },
        });

        return {
          over_limit: response.over_limit,
          total_count: response[responseKey],
        };
      } catch (err) {
        return this.setFormErrors(err);
      }
    },
  },
};
</script>
