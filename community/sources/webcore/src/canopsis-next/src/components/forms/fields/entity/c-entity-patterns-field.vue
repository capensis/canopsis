<template>
  <pattern-editor-field
    v-field="patterns"
    :disabled="disabled"
    :readonly="readonly"
    :name="name"
    :type="$constants.PATTERN_TYPES.entity"
    :required="required"
    :attributes="availableEntityAttributes"
    :with-type="withType"
    :counter="counter"
  >
    <template #append-count="">
      <v-btn
        v-if="counter && counter.count"
        text
        small
        @click="showPatternEntitiesModal"
      >
        {{ $t('common.seeEntities') }}
      </v-btn>
    </template>
  </pattern-editor-field>
</template>

<script>
import { keyBy, merge, isArray } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import {
  BASIC_ENTITY_TYPES,
  ENTITY_PATTERN_FIELDS,
  ENTITY_TYPES,
  MAX_LIMIT,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS,
  PATTERN_RULE_TYPES,
  PATTERN_STRING_OPERATORS,
} from '@/constants';

import { formGroupsToPatternRulesQuery } from '@/helpers/entities/pattern/form';

import { patternCountEntitiesModalMixin } from '@/mixins/pattern/pattern-count-entities-modal';

import PatternEditorField from '@/components/forms/fields/pattern/pattern-editor-field.vue';

const { mapActions: mapEntityCategoryActions } = createNamespacedHelpers('entityCategory');
const { mapActions: mapServiceActions } = createNamespacedHelpers('service');

export default {
  components: { PatternEditorField },
  mixins: [patternCountEntitiesModalMixin],
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Object,
      required: true,
    },
    attributes: {
      type: Array,
      default: () => [],
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
    entityTypes: {
      type: Array,
      required: false,
    },
    checkCountName: {
      type: String,
      required: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    counter: {
      type: Object,
      required: false,
    },
  },
  data() {
    return {
      categories: [],
      categoriesPending: false,
      infos: [],
    };
  },
  computed: {
    entitiesOperators() {
      return [
        PATTERN_OPERATORS.equal,
        PATTERN_OPERATORS.notEqual,
        PATTERN_OPERATORS.isOneOf,
        PATTERN_OPERATORS.isNotOneOf,
        PATTERN_OPERATORS.contains,
        PATTERN_OPERATORS.notContains,
        PATTERN_OPERATORS.regexp,
      ];
    },

    nameOptions() {
      return {
        operators: [
          ...PATTERN_STRING_OPERATORS,

          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
        ],
      };
    },

    entitiesValueField() {
      return {
        is: 'c-entity-field',
        props: {
          required: true,
          entityTypes: this.entityTypes,
        },
      };
    },

    entitiesOptions() {
      return {
        operators: this.entitiesOperators,
        defaultValue: '',
        valueField: this.entitiesValueField,
      };
    },

    componentOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
          PATTERN_OPERATORS.equal,
          PATTERN_OPERATORS.notEqual,
          PATTERN_OPERATORS.contains,
          PATTERN_OPERATORS.notContains,
          PATTERN_OPERATORS.regexp,
        ],
        defaultValue: [],
        valueField: {
          is: 'c-entity-field',
          props: {
            required: true,
            entityTypes: this.entityTypes ?? [BASIC_ENTITY_TYPES.component],
          },
        },
      };
    },

    connectorOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
          PATTERN_OPERATORS.equal,
          PATTERN_OPERATORS.notEqual,
          PATTERN_OPERATORS.contains,
          PATTERN_OPERATORS.notContains,
          PATTERN_OPERATORS.regexp,
        ],
        defaultValue: [],
        valueField: {
          is: 'c-entity-field',
          props: {
            required: true,
            entityTypes: this.entityTypes ?? [BASIC_ENTITY_TYPES.connector],
          },
        },
      };
    },

    infosOptions() {
      return {
        infos: this.infos,
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
        operators: [
          PATTERN_OPERATORS.equal,
          PATTERN_OPERATORS.notEqual,
          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
        ],
        valueField: {
          is: 'c-select-field',
          props: (rule) => {
            const isMultiple = isArray(rule?.value);

            return {
              items: this.categories,
              loading: this.categoriesPending,
              itemValue: '_id',
              itemText: 'name',
              ellipsis: true,
              multiple: isMultiple,
              deletableChips: isMultiple,
              smallChips: isMultiple,
            };
          },
        },
      };
    },

    typeOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.equal,
          PATTERN_OPERATORS.notEqual,
          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
        ],
        valueField: {
          is: 'c-entity-type-field',
          props: {
            types: Object.values(ENTITY_TYPES),
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
          options: this.nameOptions,
        },
        {
          text: this.$t('common.type'),
          value: ENTITY_PATTERN_FIELDS.type,
          options: this.typeOptions,
        },
        {
          text: this.$t('common.component'),
          value: ENTITY_PATTERN_FIELDS.component,
          options: this.componentOptions,
        },
        {
          text: this.$t('common.connector'),
          value: ENTITY_PATTERN_FIELDS.connector,
          options: this.connectorOptions,
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

    availableAttributesByValue() {
      return keyBy(this.entityAttributes, 'value');
    },

    externalAttributesByValue() {
      return keyBy(this.attributes, 'value');
    },

    availableEntityAttributes() {
      const mergedAttributes = merge(
        {},
        this.availableAttributesByValue,
        this.externalAttributesByValue,
      );

      return Object.values(mergedAttributes);
    },
  },
  mounted() {
    this.fetchCategories();
    this.fetchInfos();
  },
  methods: {
    ...mapEntityCategoryActions({ fetchCategoriesListWithoutStore: 'fetchListWithoutStore' }),
    ...mapServiceActions({ fetchEntityInfosKeysWithoutStore: 'fetchInfosKeysWithoutStore' }),

    showPatternEntitiesModal() {
      this.showEntitiesModalByPatterns({
        entity_pattern: formGroupsToPatternRulesQuery(this.patterns.groups),
      });
    },

    async fetchCategories() {
      this.categoriesPending = true;

      const { data: categories } = await this.fetchCategoriesListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.categories = categories;
      this.categoriesPending = false;
    },

    async fetchInfos() {
      const { data: infos } = await this.fetchEntityInfosKeysWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.infos = infos;
    },
  },
};
</script>
