<template>
  <pattern-editor-field
    :patterns="patterns"
    :disabled="disabled"
    :readonly="readonly"
    :name="name"
    :type="PATTERN_TYPES.pbehavior"
    :required="required"
    :attributes="pbehaviorAttributes"
    :with-type="withType"
    :alarm-counter="alarmCounter"
    :entity-counter="entityCounter"
    class="c-pbehavior-patterns-field"
    @update:patterns="$emit('update:patterns', $event)"
  >
    <template #append-count="">
      <v-btn
        v-if="alarmCounter && alarmCounter.count"
        variant="text"
        small
        @click="showPatternAlarmsModal()"
      >
        {{ t('common.seeAlarms') }}
      </v-btn>
      <v-btn
        v-if="entityCounter && entityCounter.count"
        variant="text"
        small
        @click="showPatternEntitiesModal()"
      >
        {{ t('common.seeEntities') }}
      </v-btn>
    </template>
  </pattern-editor-field>
</template>

<script>
import { computed } from 'vue';

import {
  PATTERN_OPERATORS,
  PATTERN_TYPES,
  PBEHAVIOR_PATTERN_FIELDS,
  PBEHAVIOR_TYPE_TYPES,
} from '@/constants';

import { mergePatternAttributes } from '@/helpers/entities/pattern/fields/form';

import { useI18n } from '@/hooks/i18n';

import PatternEditorField from '@/components/forms/fields/pattern/pattern-editor-field.vue';

export default {
  emits: ['update:patterns', 'show:alarms', 'show:entities'],
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
  setup(props, { emit }) {
    const { t } = useI18n();

    const nameOptions = computed(() => ({
      operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
      valueField: {
        is: 'c-pbehavior-field',
        model: 'modelValue',
        props: {
          itemTitle: 'name',
          itemValue: '_id',
          ellipsis: true,
        },
      },
    }));

    const reasonOptions = computed(() => ({
      operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
      valueField: {
        is: 'c-pbehavior-reason-field',
        model: 'modelValue',
        props: {
          ellipsis: true,
        },
      },
    }));

    const typeOptions = computed(() => ({
      operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
      valueField: {
        is: 'c-pbehavior-type-field',
        model: 'modelValue',
        props: {
          ellipsis: true,
        },
      },
    }));

    const availablePbehaviorTypeTypes = computed(() => Object.values(PBEHAVIOR_TYPE_TYPES).map(type => ({
      value: type,
      text: t(`pbehavior.types.types.${type}`),
    })));

    const canonicalTypeOptions = computed(() => ({
      operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
      valueField: {
        is: 'c-select',
        model: 'modelValue',
        props: {
          items: availablePbehaviorTypeTypes.value,
          itemTitle: 'text',
          itemValue: 'value',
        },
      },
    }));

    const defaultPbehaviorAttributes = computed(() => [
      {
        text: t('common.name'),
        value: PBEHAVIOR_PATTERN_FIELDS.id,
        options: nameOptions.value,
      },
      {
        text: t('common.reason'),
        value: PBEHAVIOR_PATTERN_FIELDS.reason,
        options: reasonOptions.value,
      },
      {
        text: t('common.type'),
        value: PBEHAVIOR_PATTERN_FIELDS.type,
        options: typeOptions.value,
      },
      {
        text: t('common.canonicalType'),
        value: PBEHAVIOR_PATTERN_FIELDS.canonicalType,
        options: canonicalTypeOptions.value,
      },
    ]);

    const pbehaviorAttributes = computed(() =>
      mergePatternAttributes(defaultPbehaviorAttributes.value, props.attributes));

    const showPatternAlarmsModal = () => {
      emit('show:alarms');
    };

    const showPatternEntitiesModal = () => {
      emit('show:entities');
    };

    return {
      t,
      PATTERN_TYPES,
      pbehaviorAttributes,
      showPatternAlarmsModal,
      showPatternEntitiesModal,
    };
  },
};
</script>
