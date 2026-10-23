<template>
  <pattern-editor-field
    :patterns="patterns"
    :disabled="disabled"
    :readonly="readonly"
    :name="name"
    :type="$constants.PATTERN_TYPES.entity"
    :required="required"
    :attributes="availableEntityAttributes"
    :with-type="withType"
    :alarm-counter="alarmCounter"
    :entity-counter="entityCounter"
    class="c-entity-patterns-field"
    @update:patterns="$emit('update:patterns', $event)"
  >
    <template #append-count="">
      <v-btn
        v-if="alarmCounter && alarmCounter.count"
        variant="text"
        small
        @click="showPatternAlarmsModal()"
      >
        {{ $t('common.seeAlarms') }}
      </v-btn>
      <v-btn
        v-if="entityCounter && entityCounter.count"
        variant="text"
        small
        @click="showPatternEntitiesModal()"
      >
        {{ $t('common.seeEntities') }}
      </v-btn>
    </template>
  </pattern-editor-field>
</template>

<script>
import { isArray } from 'lodash';

import {
  BASIC_ENTITY_TYPES,
  ENTITY_PATTERN_FIELDS,
  ENTITY_TYPES,
  MAX_LIMIT,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS,
  PATTERN_EXISTS_OPERATORS,
  PATTERN_RULE_TYPES,
  PATTERN_STRING_OPERATORS,
  ADVANCED_SEARCH_INFOS_TYPES_TO_PATTERNS_FIELD_TYPES,
  ENTITY_PATTERN_FIELD_TYPES,
} from '@/constants';

import { getMapEntityText } from '@/helpers/entities/map/list';
import { mergePatternAttributes } from '@/helpers/entities/pattern/fields/form';

import { useEntityCategory } from '@/hooks/store/modules/entity-category';
import { useService } from '@/hooks/store/modules/service';
import { entitiesEntityInfoPropertyMixin } from '@/mixins/entities/entity-info-property';

import PatternEditorField from '@/components/forms/fields/pattern/pattern-editor-field.vue';

export default {
  emits: ['update:patterns', 'show:alarms', 'show:entities'],
  components: { PatternEditorField },
  mixins: [entitiesEntityInfoPropertyMixin],
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
    alarmCounter: {
      type: Object,
      required: false,
    },
    entityCounter: {
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
        PATTERN_OPERATORS.beginWith,
        PATTERN_OPERATORS.notBeginWith,
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
          itemText: getMapEntityText,
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
          PATTERN_OPERATORS.beginWith,
          PATTERN_OPERATORS.notBeginWith,
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
          PATTERN_OPERATORS.beginWith,
          PATTERN_OPERATORS.notBeginWith,
          PATTERN_OPERATORS.regexp,
        ],
        defaultValue: [],
        valueField: {
          is: 'c-entity-field',
          props: {
            required: true,
            entityTypes: [BASIC_ENTITY_TYPES.connector],
          },
        },
      };
    },

    infosWithDefinedTypes() {
      return this.infos.map(({ type, ...info }) => ({
        ...info,
        definedType: ADVANCED_SEARCH_INFOS_TYPES_TO_PATTERNS_FIELD_TYPES[type],
      }));
    },

    infosOptions() {
      return {
        infos: this.infosWithDefinedTypes,
        type: PATTERN_RULE_TYPES.infos,
        inputTypes: ENTITY_PATTERN_FIELD_TYPES,
      };
    },

    componentInfosOptions() {
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
          is: 'c-select',
          props: (rule) => {
            const isMultiple = isArray(rule?.value);

            return {
              items: this.categories,
              loading: this.categoriesPending,
              itemValue: '_id',
              itemTitle: 'name',
              ellipsis: !isMultiple,
              chips: isMultiple,
              multiple: isMultiple,
              closableChips: isMultiple,
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

    aliasesAttributes() {
      return this.entityInfoPropertiesWithAlias.map(item => ({
        text: item.alias,
        value: item.alias,
        alias: true,
        originalValue: item.name,
        definedType: ADVANCED_SEARCH_INFOS_TYPES_TO_PATTERNS_FIELD_TYPES[item.type],
        options: {
          inputTypes: ENTITY_PATTERN_FIELD_TYPES,
        },
      }));
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
          text: this.$t('common.componentInfo', 2),
          value: ENTITY_PATTERN_FIELDS.componentInfos,
          options: this.componentInfosOptions,
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
          text: this.$t('entity.fields.pause'),
          value: ENTITY_PATTERN_FIELDS.pause,
          options: { operators: PATTERN_EXISTS_OPERATORS },
        },
        {
          text: this.$t('entity.fields.pauseAuthor'),
          value: ENTITY_PATTERN_FIELDS.pauseAuthor,
        },
        {
          text: this.$t('entity.fields.pauseInitiator'),
          value: ENTITY_PATTERN_FIELDS.pauseInitiator,
        },
        {
          text: this.$t('common.lastEventDate'),
          value: ENTITY_PATTERN_FIELDS.lastEventDate,
          options: this.dateOptions,
        },
        ...this.aliasesAttributes,
      ];
    },

    availableEntityAttributes() {
      return mergePatternAttributes(this.entityAttributes, this.attributes);
    },
  },
  mounted() {
    this.fetchCategories();
    this.fetchInfos();
    this.fetchAllEntityInfoPropertiesList();
  },
  methods: {
    showPatternAlarmsModal() {
      this.$emit('show:alarms');
    },

    showPatternEntitiesModal() {
      this.$emit('show:entities');
    },

    async fetchCategories() {
      const { fetchCategoriesListWithoutStore } = useEntityCategory();

      this.categoriesPending = true;

      const { data: categories } = await fetchCategoriesListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.categories = categories;
      this.categoriesPending = false;
    },

    async fetchInfos() {
      const { fetchEntityInfosKeysWithoutStore } = useService();
      const { data: infos } = await fetchEntityInfosKeysWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.infos = infos;
    },
  },
};
</script>
