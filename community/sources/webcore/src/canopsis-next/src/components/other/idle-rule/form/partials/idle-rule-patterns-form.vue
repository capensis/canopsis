<template lang="pug">
  c-patterns-field.mt-2(
    v-field="form",
    :with-alarm="!isEntityType",
    some-required,
    with-entity
  )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { PATTERNS_TYPES } from '@/constants';

import { formValidationHeaderMixin, validationErrorsMixinCreator } from '@/mixins/form';

const { mapActions } = createNamespacedHelpers('idleRules');

export default {
  provide() {
    return {
      $checkEntitiesCountForPatternsByType: this.checkEntitiesCountForPatternsByType,
    };
  },
  inject: ['$validator'],
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
