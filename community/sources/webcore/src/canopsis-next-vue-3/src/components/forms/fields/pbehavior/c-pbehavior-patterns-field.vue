<template>
  <pattern-editor-field
    :patterns="patterns"
    :disabled="disabled"
    :readonly="readonly"
    :name="name"
    :type="$constants.PATTERN_TYPES.pbehavior"
    :required="required"
    :attributes="pbehaviorAttributes"
    :with-type="withType"
    :alarm-counter="alarmCounter"
    :entity-counter="entityCounter"
    class="c-pbehavior-patterns-field"
    @update:patterns="$emit('update:value', $event)"
  >
    <template #append-count="">
      <v-btn
        v-if="alarmCounter && alarmCounter.count"
        text
        small
        @click="showPatternAlarmsModal()"
      >
        {{ $t('common.seeAlarms') }}
      </v-btn>
      <v-btn
        v-if="entityCounter && entityCounter.count"
        text
        small
        @click="showPatternEntitiesModal()"
      >
        {{ $t('common.seeEntities') }}
      </v-btn>
    </template>
  </pattern-editor-field>
</template>

<script>

import { MAX_LIMIT, PATTERN_OPERATORS, PBEHAVIOR_PATTERN_FIELDS, PBEHAVIOR_TYPE_TYPES } from '@/constants';
import { usePbehavior } from '@/hooks/store/modules/pbehavior';

import PatternEditorField from '@/components/forms/fields/pattern/pattern-editor-field.vue';

export default {
  emits: ['update:patterns', 'update:value'],
  components: { PatternEditorField },
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
    attributes: {
      type: Array,
      default: () => [],
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
          is: 'c-select',
          props: {
            items: this.pbehaviors,
            loading: this.pbehaviorsPending,
            itemValue: '_id',
            itemTitle: 'name',
            ellipsis: true,
          },
        },
      };
    },

    reasonOptions() {
      return {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-select',
          props: {
            items: this.pbehaviorReasons,
            loading: this.pbehaviorReasonsPending,
            itemValue: '_id',
            itemTitle: 'name',
            ellipsis: true,
          },
        },
      };
    },

    typeOptions() {
      return {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-select',
          props: {
            items: this.pbehaviorTypes,
            loading: this.pbehaviorTypesPending,
            itemValue: '_id',
            itemTitle: 'name',
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
          is: 'c-select',
          props: {
            items: this.availablePbehaviorTypeTypes,
          },
        },
      };
    },

    pbehaviorAttributes() {
      return [
        ...this.attributes,
        {
          text: this.$t('common.name'),
          value: PBEHAVIOR_PATTERN_FIELDS.id,
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
    showPatternAlarmsModal() {
      this.$emit('show:alarms');
    },

    showPatternEntitiesModal() {
      this.$emit('show:entities');
    },

    async fetchPbehaviors() {
      const { fetchPbehaviorsListWithoutStore } = usePbehavior();

      this.pbehaviorsPending = true;

      const { data: pbehaviors } = await fetchPbehaviorsListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.pbehaviors = pbehaviors;
      this.pbehaviorsPending = false;
    },

    async fetchPbehaviorReasons() {
      const { fetchPbehaviorReasonsListWithoutStore } = usePbehavior();

      this.pbehaviorReasonsPending = true;

      const { data: reasons } = await fetchPbehaviorReasonsListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.pbehaviorReasons = reasons;
      this.pbehaviorReasonsPending = false;
    },

    async fetchPbehaviorTypes() {
      const { fetchPbehaviorTypesListWithoutStore } = usePbehavior();

      this.pbehaviorTypesPending = true;

      const { data: types } = await fetchPbehaviorTypesListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.pbehaviorTypes = types;
      this.pbehaviorTypesPending = false;
    },
  },
};
</script>
