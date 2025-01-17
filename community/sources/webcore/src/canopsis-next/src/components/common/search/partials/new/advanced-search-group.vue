<template>
  <v-layout
    :class="{ 'c-new-advanced-search__group--union': union }"
    class="c-new-advanced-search__group"
  >
    <advanced-search-chip
      v-for="chip in chips"
      v-bind="chip.bind"
      :key="chip.key"
      v-on="chip.on"
    />
  </v-layout>
</template>

<script>
import { keyBy, uniq } from 'lodash';
import { computed, ref, watch, nextTick, onMounted, onBeforeUnmount, inject } from 'vue';
import { Validator } from 'vee-validate';

import {
  ADVANCED_SEARCH_UNION_CONDITIONS,
  ALARM_FIELDS, ALARM_PATTERN_FIELDS,
  PATTERN_DURATION_OPERATORS,
  PATTERN_FIELD_TYPES,
  PATTERN_OPERATORS,
  PATTERN_QUICK_RANGES,
} from '@/constants';

import AdvancedSearchChip from './advanced-search-chip.vue';
import {
  isArrayCondition,
  isDatePatternRuleField, isDurationPatternRuleField,
  isInfosPatternRuleField,
  isValueInfosPatternRuleField
} from '@/helpers/entities/pattern/form';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';
import { uid } from '@/helpers/uid';


const patterns = {
  groups: [
    {
      key: 'group1',
      rules: [
        {
          filled: ['attribute', 'operation', 'value'],
          key: 'd42ace4e0d3',
          attribute: '',
          dictionary: '',
          duration: {
            unit: 's',
            value: 1,
          },
          field: '',
          fieldType: 'string',
          operator: '',
          range: {
            from: 0,
            to: 0,
            type: 'last1Hour',
          },
          value: '',
        },
      ],
    },
  ],
};

const ADVANCED_SEARCH_CHIP_TYPES = {
  attribute: 'attribute',
  dictionary: 'dictionary',
  operator: 'operator',
  fieldType: 'fieldType',
  value: 'value',
  duration: 'duration',
  range: 'range',
  union: 'union',
};

export const getNextType = (attribute, type) => {
  if (type === ADVANCED_SEARCH_CHIP_TYPES.union) {
    return null;
  }

  if (!attribute) {
    return ADVANCED_SEARCH_CHIP_TYPES.attribute;
  }

  switch (type) {
    case ADVANCED_SEARCH_CHIP_TYPES.attribute:
      if (isDatePatternRuleField(attribute)) {
        return ADVANCED_SEARCH_CHIP_TYPES.range;
      }

      if (isValueInfosPatternRuleField(attribute)) {
        return ADVANCED_SEARCH_CHIP_TYPES.fieldType;
      }

      if (attribute === ALARM_PATTERN_FIELDS.ticketData) {
        return ADVANCED_SEARCH_CHIP_TYPES.dictionary;
      }

      return ADVANCED_SEARCH_CHIP_TYPES.operator;

    case ADVANCED_SEARCH_CHIP_TYPES.dictionary:
    case ADVANCED_SEARCH_CHIP_TYPES.fieldType:
      return ADVANCED_SEARCH_CHIP_TYPES.operator;

    case ADVANCED_SEARCH_CHIP_TYPES.operator:
      if (isDurationPatternRuleField(attribute)) {
        return ADVANCED_SEARCH_CHIP_TYPES.duration;
      }

      return ADVANCED_SEARCH_CHIP_TYPES.value;
    default:
      return null;
  }
};

export default {
  components: { AdvancedSearchChip },
  model: {
    prop: 'rule',
    event: 'input',
  },
  props: {
    rule: {
      type: Object,
      default: () => ({}),
    },
    attributes: {
      type: Array,
      default: () => [],
    },
    activeKey: {
      type: String,
      default: '',
    },
    union: {
      type: Boolean,
      default: false,
    },
    inputTypes: {
      type: Array,
      default: () => [
        { value: PATTERN_FIELD_TYPES.string },
        { value: PATTERN_FIELD_TYPES.number },
        { value: PATTERN_FIELD_TYPES.boolean },
        { value: PATTERN_FIELD_TYPES.stringArray },
      ],
    },
    intervalRanges: {
      type: Array,
      default: () => PATTERN_QUICK_RANGES,
    },
    allowOr: {
      type: Boolean,
      default: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const validator = inject('$validator', new Validator());

    const { t } = useI18n();
    const { updateField, updateModel } = useModelField(props, emit);
    const activeType = ref(props.union ? ADVANCED_SEARCH_CHIP_TYPES.union : ADVANCED_SEARCH_CHIP_TYPES.attribute);

    const nextType = computed(() => getNextType(props.rule.attribute, activeType.value));
    const attributesMap = computed(() => keyBy(props.attributes, 'value'));
    const isFinishedGroup = computed(() => !activeType.value);
    const preparedOperators = computed(() => (
      (attributesMap.value[props.rule.attribute]?.operators ?? []).map(operator => ({
        text: t(`common.operators.${operator}`),
        value: operator,
      }))
    ));

    const preparedIntervalRanges = computed(() => (
      props.intervalRanges.map(range => ({
        ...range,
        text: t(`quickRanges.types.${range.value}`),
      }))
    ));

    const preparedInputTypes = computed(() => (
      props.inputTypes.map(type => ({
        ...type,
        text: t(`common.mixedField.${type.value}`),
      }))
    ));

    const preparedUnionItems = computed(() => (
      Object.values(ADVANCED_SEARCH_UNION_CONDITIONS).map(value => ({
        value,
        text: value,
      }))
    ));

    const itemsByType = computed(() => ({
      [ADVANCED_SEARCH_CHIP_TYPES.attribute]: props.attributes,
      [ADVANCED_SEARCH_CHIP_TYPES.operator]: preparedOperators.value,
      [ADVANCED_SEARCH_CHIP_TYPES.range]: preparedIntervalRanges.value,
      [ADVANCED_SEARCH_CHIP_TYPES.fieldType]: preparedInputTypes.value,
      [ADVANCED_SEARCH_CHIP_TYPES.union]: preparedUnionItems.value,
    }));

    const goToNextType = (steps = 1) => {
      for (let i = 0; i < steps; i += 1) {
        activeType.value = getNextType(props.rule.attribute, activeType.value);
      }
    };

    const selectChipItem = (value, type) => updateField(type, value);

    const selectItem = (value) => {
      const filled = [...(props.rule.filled ?? []), activeType.value];
      const preparedRule = { ...props.rule };
      let actionTypeSteps = 1;

      if (activeType.value === ADVANCED_SEARCH_CHIP_TYPES.operator && isArrayCondition(value)) {
        filled.push(ADVANCED_SEARCH_CHIP_TYPES.value);
        preparedRule[ADVANCED_SEARCH_CHIP_TYPES.value] = [];
        actionTypeSteps = 2;
      }

      updateModel({
        ...preparedRule,
        filled: uniq(filled),
        [activeType.value]: value,
      });


      if (!isFinishedGroup.value) {
        nextTick(() => goToNextType(actionTypeSteps));
      }
    };

    const clickChip = key => emit('click:chip', key);

    const remove = () => console.log('REMOVE') || emit('remove');

    const chips = computed(() => {
      const result = (props.rule.filled ?? []).map((type, index, filled) => {
        const key = `${props.rule.key}.${type}`;

        return {
          key,
          bind: {
            value: props.rule[type],
            active: key === props.activeKey,
            items: itemsByType.value[type],
            allowText: !itemsByType.value[type]?.length,
            closable: index === filled.length - 1 && !props.union,
            disabled: props.disabled,
            multiple: type === ADVANCED_SEARCH_CHIP_TYPES.value && isArrayCondition(props.rule.operator),
            color: validator.errors.has(props.rule.key) ? 'error' : undefined,
          },
          on: {
            input: value => selectChipItem(value, type),
            click: () => clickChip(key),
            focusout: () => emit('focusout'),
            next: () => setTimeout(() => emit('next', isFinishedGroup.value)), // TODO: refactor
            close: remove,
          },
        };
      });

      if (!isFinishedGroup.value) {
        result.push({
          key: uid(),
          bind: {
            active: true,
            alwaysActive: true,
            items: itemsByType.value[activeType.value],
            allowText: !itemsByType.value[activeType.value]?.length,
          },
          on: {
            input: selectItem,
            next: () => setTimeout(() => emit('next', isFinishedGroup.value)), // TODO: refactor
          },
        });
      }

      return result;
    });

    watch(() => props.union, (union) => {
      activeType.value = union ? ADVANCED_SEARCH_CHIP_TYPES.union : ADVANCED_SEARCH_CHIP_TYPES.attribute;
    });

    onMounted(() => {
      validator.attach({
        name: props.rule.key,
        rules: 'required:true',
        getter: () => isFinishedGroup.value,
      });
    });

    onBeforeUnmount(() => {
      validator.detach(props.rule.key);
    });

    return {
      activeType,
      nextType,
      isFinishedGroup,
      chips,
      remove,
      selectChipItem,
      selectItem,
      clickChip,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-new-advanced-search__group {
  &:hover ::v-deep {
    .theme--light.v-chip:before {
      opacity: 0.04;
    }
    .theme--dark.v-chip:before {
      opacity: 0.08;
    }
  }

  &--union ::v-deep {
    .v-chip {
      background: var(--v-application-background-base) !important;

      &.theme--light {
        border: 1px solid var(--v-text-light-primary, rgba(0, 0, 0, 0.87));
      }

      &.theme--dark {
        border: 1px solid var(--v-text-dark-primary, #FFFFFF);
      }
    }
  }
}
</style>
