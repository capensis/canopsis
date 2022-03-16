<template lang="pug">
  c-pattern-groups-field(
    v-field="groups",
    :disabled="disabled",
    :name="name",
    :attributes="entityAttributes"
  )
</template>

<script>
import {
  ENTITY_PATTERN_FIELDS,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS,
  PATTERN_RULE_TYPES,
} from '@/constants';

export default {
  model: {
    prop: 'groups',
    event: 'input',
  },
  props: {
    groups: {
      type: Array,
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
        props: {
          // TODO: Should be replaced on API data
          items: [
            { name: 'Entity 1', _id: 'entity_1' },
            { name: 'Entity 2', _id: 'entity_2' },
            { name: 'Entity 3', _id: 'entity_3' },
            { name: 'Entity 4', _id: 'entity_4' },
          ],
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
        valueField: {
          is: 'c-entity-category-field',
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
};
</script>
