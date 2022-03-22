<template lang="pug">
  c-patterns-editor-field(
    v-field="patterns",
    :disabled="disabled",
    :name="name",
    :type="$constants.PATTERN_TYPES.entity",
    :required="required",
    :attributes="entityAttributes",
    :with-type="withType"
  )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import {
  ENTITY_PATTERN_FIELDS,
  MAX_LIMIT,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS,
  PATTERN_RULE_TYPES,
} from '@/constants';

const { mapActions: entityCategoryMapActions } = createNamespacedHelpers('entityCategory');

export default {
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Object,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      required: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    withType: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      categories: [],
      categoriesPending: false,
    };
  },
  computed: {
    entitiesOperators() {
      return [
        PATTERN_OPERATORS.equal,
        PATTERN_OPERATORS.notEqual,
        PATTERN_OPERATORS.hasOneOf,
        PATTERN_OPERATORS.hasNot,
      ];
    },

    entitiesValueField() {
      return {
        is: 'c-entity-field',
      };
    },

    entitiesOptions() {
      return {
        operators: this.entitiesOperators,
        defaultValue: '',
        valueField: this.entitiesValueField,
      };
    },

    infosOptions() {
      return {
        // TODO: Should be replaced on API data
        infos: ['infos 1', 'infos 2'],
        type: PATTERN_RULE_TYPES.infos,
      };
    },

    dateOptions() {
      return {
        type: PATTERN_RULE_TYPES.date,
      };
    },

    impactLevelOptions() {
      return {
        type: PATTERN_RULE_TYPES.string,
        operators: PATTERN_NUMBER_OPERATORS,
        valueField: {
          is: 'c-impact-level-field',
          required: true,
        },
      };
    },

    categoryOptions() {
      return {
        type: PATTERN_RULE_TYPES.string,
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-select-field',
          props: {
            items: this.categories,
            loading: this.categoriesPending,
            itemValue: '_id',
            itemText: 'name',
            ellipsis: true,
          },
        },
      };
    },

    entityAttributes() {
      return [
        {
          text: this.$t('common.id'),
          value: ENTITY_PATTERN_FIELDS.id,
          options: this.entitiesOptions,
        },
        {
          text: this.$t('common.name'),
          value: ENTITY_PATTERN_FIELDS.name,
        },
        {
          text: this.$tc('common.impact', 2),
          value: ENTITY_PATTERN_FIELDS.impact,
          options: this.entitiesOptions,
        },
        {
          text: this.$tc('common.depend', 2),
          value: ENTITY_PATTERN_FIELDS.depends,
          options: this.entitiesOptions,
        },
        {
          text: this.$t('common.infos'),
          value: ENTITY_PATTERN_FIELDS.infos,
          options: this.infosOptions,
        },
        {
          text: this.$tc('common.componentInfo', 2),
          value: ENTITY_PATTERN_FIELDS.componentInfos,
          options: this.infosOptions,
        },
        {
          text: this.$t('common.category'),
          value: ENTITY_PATTERN_FIELDS.category,
          options: this.categoryOptions,
        },
        {
          text: this.$t('common.impactLevel'),
          value: ENTITY_PATTERN_FIELDS.impactLevel,
          options: this.impactLevelOptions,
        },
        {
          text: this.$t('common.lastEventDate'),
          value: ENTITY_PATTERN_FIELDS.lastEventDate,
          options: this.dateOptions,
        },
      ];
    },
  },
  mounted() {
    this.fetchCategories();
  },
  methods: {
    ...entityCategoryMapActions({ fetchCategoriesListWithoutStore: 'fetchListWithoutStore' }),

    async fetchCategories() {
      this.categoriesPending = true;

      const { data: categories } = await this.fetchCategoriesListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.categories = categories;
      this.categoriesPending = false;
    },
  },
};
</script>
