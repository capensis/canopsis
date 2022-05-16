<template lang="pug">
  c-patterns-field(
    v-field="form",
    :with-alarm="!isEntityType",
    :alarm-attributes="alarmAttributes",
    :entity-attributes="entityAttributes",
    some-required,
    with-entity
  )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { ALARM_PATTERN_FIELDS, ENTITY_PATTERN_FIELDS, PATTERNS_TYPES, QUICK_RANGES } from '@/constants';

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
  computed: {
    alarmAttributes() {
      return [
        {
          value: ALARM_PATTERN_FIELDS.creationDate,
          options: {
            intervalRanges: [QUICK_RANGES.custom],
          },
        },
        {
          value: ALARM_PATTERN_FIELDS.lastUpdateDate,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.lastEventDate,
          options: { disabled: true },
        },
        {
          value: ALARM_PATTERN_FIELDS.resolvedAt,
          options: { disabled: true },
        },
      ];
    },

    entityAttributes() {
      return [
        {
          value: ENTITY_PATTERN_FIELDS.lastEventDate,
          options: {
            disabled: true,
          },
        },
      ];
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
