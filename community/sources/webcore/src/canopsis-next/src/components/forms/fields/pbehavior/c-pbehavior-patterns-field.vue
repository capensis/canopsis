<template>
  <pattern-editor-field
    v-field="patterns"
    :disabled="disabled"
    :readonly="readonly"
    :name="name"
    :type="$constants.PATTERN_TYPES.pbehavior"
    :required="required"
    :attributes="pbehaviorAttributes"
    :with-type="withType"
    :counter="counter"
  />
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT, PATTERN_OPERATORS, PBEHAVIOR_PATTERN_FIELDS, PBEHAVIOR_TYPE_TYPES } from '@/constants';

import PatternEditorField from '@/components/forms/fields/pattern/pattern-editor-field.vue';

const { mapActions: pbehaviorMapActions } = createNamespacedHelpers('pbehavior');
const { mapActions: pbehaviorReasonMapActions } = createNamespacedHelpers('pbehaviorReasons');
const { mapActions: pbehaviorTypeMapActions } = createNamespacedHelpers('pbehaviorTypes');

export default {
  components: { PatternEditorField },
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
      pbehaviors: [],
      pbehaviorsPending: false,

      pbehaviorReasons: [],
      pbehaviorReasonsPending: false,

      pbehaviorTypes: [],
      pbehaviorTypesPending: false,
    };
  },
  computed: {
    nameOptions() {
      return {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-select-field',
          props: {
            items: this.pbehaviors,
            loading: this.pbehaviorsPending,
            itemValue: '_id',
            itemText: 'name',
            ellipsis: true,
          },
        },
      };
    },

    reasonOptions() {
      return {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-select-field',
          props: {
            items: this.pbehaviorReasons,
            loading: this.pbehaviorReasonsPending,
            itemValue: '_id',
            itemText: 'name',
            ellipsis: true,
          },
        },
      };
    },

    typeOptions() {
      return {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-select-field',
          props: {
            items: this.pbehaviorTypes,
            loading: this.pbehaviorTypesPending,
            itemValue: '_id',
            itemText: 'name',
            ellipsis: true,
          },
        },
      };
    },

    availablePbehaviorTypeTypes() {
      return Object.values(PBEHAVIOR_TYPE_TYPES).map(type => ({
        value: type,
        text: this.$t(`pbehavior.types.types.${type}`),
      }));
    },

    canonicalTypeOptions() {
      return {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-select-field',
          props: {
            items: this.availablePbehaviorTypeTypes,
          },
        },
      };
    },

    pbehaviorAttributes() {
      return [
        {
          text: this.$t('common.name'),
          value: PBEHAVIOR_PATTERN_FIELDS.name,
          options: this.nameOptions,
        },
        {
          text: this.$t('common.reason'),
          value: PBEHAVIOR_PATTERN_FIELDS.reason,
          options: this.reasonOptions,
        },
        {
          text: this.$t('common.type'),
          value: PBEHAVIOR_PATTERN_FIELDS.type,
          options: this.typeOptions,
        },
        {
          text: this.$t('common.canonicalType'),
          value: PBEHAVIOR_PATTERN_FIELDS.canonicalType,
          options: this.canonicalTypeOptions,
        },
      ];
    },
  },
  mounted() {
    this.fetchPbehaviors();
    this.fetchPbehaviorReasons();
    this.fetchPbehaviorTypes();
  },
  methods: {
    ...pbehaviorMapActions({ fetchPbehaviorsListWithoutStore: 'fetchListWithoutStore' }),
    ...pbehaviorReasonMapActions({ fetchPbehaviorReasonsListWithoutStore: 'fetchListWithoutStore' }),
    ...pbehaviorTypeMapActions({ fetchPbehaviorTypesListWithoutStore: 'fetchListWithoutStore' }),

    async fetchPbehaviors() {
      this.pbehaviorsPending = true;

      const { data: pbehaviors } = await this.fetchPbehaviorsListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.pbehaviors = pbehaviors;
      this.pbehaviorsPending = false;
    },

    async fetchPbehaviorReasons() {
      this.pbehaviorReasonsPending = true;

      const { data: reasons } = await this.fetchPbehaviorReasonsListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.pbehaviorReasons = reasons;
      this.pbehaviorReasonsPending = false;
    },

    async fetchPbehaviorTypes() {
      this.pbehaviorTypesPending = true;

      const { data: types } = await this.fetchPbehaviorTypesListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.pbehaviorTypes = types;
      this.pbehaviorTypesPending = false;
    },
  },
};
</script>
