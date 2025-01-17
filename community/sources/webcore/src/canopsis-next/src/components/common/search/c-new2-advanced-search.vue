<template>
  <div
    :class="[wrapperClass, themeClasses]"
    class="c-new-advanced-search v-input v-input--hide-details theme--light
    v-text-field v-text-field--single-line v-text-field--is-booted v-select v-autocomplete primary--text"
  >
    <div class="v-input__control">
      <div class="v-input__slot">
        <v-layout align-center justify-start wrap>
          <div class="gap-1">
            <advanced-search-group
              v-for="(group, index) in groups"
              :key="group.key"
              :value="group"
              :fields="fields"
              :active-key="activeKey"
              @input="update($event, index)"
              @mousedown="mousedownChip"
            />
          </div>
        </v-layout>
      </div>
    </div>
  </div>
</template>

<script>
import { computed, ref, set } from 'vue';
import Themeable from 'vuetify/lib/mixins/themeable';

import {
  ALARM_EVENT_INITIATORS,
  ALARM_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ALARM_GROUPED_ADVANCED_SEARCH_FIELDS,
  PATTERN_DURATION_OPERATORS,
  PATTERN_OPERATORS,
  PATTERN_QUICK_RANGES,
  PATTERN_STRING_OPERATORS,
  QUICK_RANGES,
} from '@/constants';

import { uid } from '@/helpers/uid';

import { useI18n } from '@/hooks/i18n';

import AdvancedSearchChip from 'src/components/common/search/partials/new/chips/advanced-search-chip.vue';
import AdvancedSearchMultipleChip from './partials/new/advanced-search-multiple-chip.vue';
import AdvancedSearchGroup from '@/components/common/search/partials/new/advanced-search-group.vue';

const ADVANCED_SEARCH_STEP_TYPES = {
  field: 'field',
  operator: 'operator',
  valueType: 'valueType',
  value: 'value',
  valueSecond: 'valueSecond',
  union: 'union',
};

const generateUntypicalDateTypes = () => [
  ALARM_FIELDS.creationDate,
  ALARM_FIELDS.lastUpdateDate,
  ALARM_FIELDS.lastEventDate,
  ALARM_FIELDS.ackAt,
  ALARM_FIELDS.resolved,
  ALARM_FIELDS.activationDate,
].reduce((acc, field) => {
  [
    QUICK_RANGES.last15Minutes,
    QUICK_RANGES.last30Minutes,
    QUICK_RANGES.last1Hour,
    QUICK_RANGES.last3Hour,
    QUICK_RANGES.last3Hour,
    QUICK_RANGES.last6Hour,
    QUICK_RANGES.last12Hour,
    QUICK_RANGES.last24Hour,
    QUICK_RANGES.last2Days,
    QUICK_RANGES.last7Days,
    QUICK_RANGES.last30Days,
    QUICK_RANGES.last1Year,
  ].forEach(({ value }) => {
    acc[`${field}.${value}`] = ADVANCED_SEARCH_STEP_TYPES.union;

    return acc;
  });

  acc[QUICK_RANGES.custom.value] = ADVANCED_SEARCH_STEP_TYPES.valueSecond;

  return acc;
});

const ARRAY_OPERATORS = [
  PATTERN_OPERATORS.isOneOf,
  PATTERN_OPERATORS.isNotOneOf,
];

export const useAlarmAdvancedSearchVariables = () => {
  const STRING_WITH_ONE_OF_OPERATORS = [
    ...PATTERN_STRING_OPERATORS,

    PATTERN_OPERATORS.isOneOf,
    PATTERN_OPERATORS.isNotOneOf,
  ];
  const STRING_WITH_EXIST_AND_ONE_OF_OPERATORS = [
    ...STRING_WITH_ONE_OF_OPERATORS,

    PATTERN_OPERATORS.exist,
  ];

  const NUMBER_OPERATORS = [
    PATTERN_OPERATORS.equal,
    PATTERN_OPERATORS.notEqual,
    PATTERN_OPERATORS.higher,
    PATTERN_OPERATORS.lower,
  ];

  const TICKET_OPERATORS = [
    PATTERN_OPERATORS.ticketAssociated,
    PATTERN_OPERATORS.ticketNotAssociated,
  ];

  const ACK_OPTIONS = [
    PATTERN_OPERATORS.ticketAssociated,
    PATTERN_OPERATORS.ticketNotAssociated,
  ];

  const USER_OPERATORS = [
    PATTERN_OPERATORS.equal,
    PATTERN_OPERATORS.notEqual,
    PATTERN_OPERATORS.isOneOf,
    PATTERN_OPERATORS.isNotOneOf,
  ];

  const ENTITY_OPERATORS = [
    ...PATTERN_STRING_OPERATORS,

    PATTERN_OPERATORS.isOneOf,
    PATTERN_OPERATORS.isNotOneOf,
  ];

  const DATE_OPTIONS = {
    operators: PATTERN_QUICK_RANGES,
  };

  const INITIATOR_OPTIONS = {
    operators: USER_OPERATORS,
    items: Object.values(ALARM_EVENT_INITIATORS).map(initiator => ({ value: initiator, text: initiator })),
  };

  const ADVANCED_SEARCH_OPTIONS_MAP = {
    /**
     * Basic
     */
    [ALARM_FIELDS.displayName]: {
      operators: STRING_WITH_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.connector]: {},
    [ALARM_FIELDS.connectorName]: {},
    [ALARM_FIELDS.component]: {},
    [ALARM_FIELDS.resource]: {},
    [ALARM_FIELDS.state]: {},
    [ALARM_FIELDS.status]: {},
    [ALARM_FIELDS.tags]: {},
    [ALARM_FIELDS.infos]: {},
    [ALARM_FIELDS.meta]: {},
    [ALARM_FIELDS.changeState]: {},
    [ALARM_FIELDS.totalStateChanges]: {},

    /**
     * Messages
     */
    [ALARM_FIELDS.output]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.longOutput]: {
      operators: STRING_WITH_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.initialOutput]: {
      operators: STRING_WITH_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.initialLongOutput]: {
      operators: STRING_WITH_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.lastComment]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.lastCommentInitiator]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },

    /**
     * Ticket
     */
    [ALARM_FIELDS.ticketMessage]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.ticketInitiator]: INITIATOR_OPTIONS,
    [ALARM_FIELDS.ticketValue]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.ticket]: {
      operators: [
        PATTERN_OPERATORS.ticketAssociated,
        PATTERN_OPERATORS.ticketNotAssociated,
      ],
    },

    /**
     * Ticket
     */
    [ALARM_FIELDS.creationDate]: {
      operators: PATTERN_QUICK_RANGES,
    },
    [ALARM_FIELDS.lastUpdateDate]: {
      operators: PATTERN_QUICK_RANGES,
    },
    [ALARM_FIELDS.lastEventDate]: {
      operators: PATTERN_QUICK_RANGES,
    },
    [ALARM_FIELDS.ackAt]: {
      operators: PATTERN_QUICK_RANGES,
    },
    [ALARM_FIELDS.resolved]: {
      operators: PATTERN_QUICK_RANGES,
    },
    [ALARM_FIELDS.activationDate]: {
      operators: PATTERN_QUICK_RANGES,
    },
    [ALARM_FIELDS.duration]: {
      operators: PATTERN_DURATION_OPERATORS,
    },

    /**
     * Actions
     */
    [ALARM_FIELDS.ack]: {
      operators: [
        PATTERN_OPERATORS.acked,
        PATTERN_OPERATORS.notAcked,
      ],
    },
    [ALARM_FIELDS.ackBy]: {
      operators: USER_OPERATORS,
      valueFetch: () => ({ data: [{ _id: 'root' }], meta: { page_count: 1 } }),
    },
    [ALARM_FIELDS.ackMessage]: {
      operators: STRING_WITH_EXIST_AND_ONE_OF_OPERATORS,
    },
    [ALARM_FIELDS.ackInitiator]: INITIATOR_OPTIONS,
    [ALARM_FIELDS.canceled]: {
      operators: [
        PATTERN_OPERATORS.canceled,
        PATTERN_OPERATORS.notCanceled,
      ],
    },
    [ALARM_FIELDS.canceledInitiator]: INITIATOR_OPTIONS,
    [ALARM_FIELDS.activated]: {
      operators: [
        PATTERN_OPERATORS.activated,
        PATTERN_OPERATORS.inactive,
      ],
    },
    [ALARM_FIELDS.snooze]: {
      operators: [
        PATTERN_OPERATORS.snoozed,
        PATTERN_OPERATORS.notSnoozed,
      ],
    },

    /**
     * Entity
     */
    [ALARM_FIELDS.entityId]: {},
    [ALARM_FIELDS.entityName]: {},
    [ALARM_FIELDS.entityCategoryName]: {},
    [ALARM_FIELDS.entityType]: {},
    [ALARM_FIELDS.entityComponent]: {},
    [ALARM_FIELDS.entityConnector]: {},
    [ALARM_FIELDS.entityImpactLevel]: {},
    [ALARM_FIELDS.entityInfos]: {
      items: [
        {
          value: 'Something',
          text: 'Something',
          items: [{ value: 'value', text: 'Value' }] },
      ],
    },
    [ALARM_FIELDS.entityComponentInfos]: {},

    /**
     * Pbehavior
     */
    [ALARM_FIELDS.pbehaviorInfoName]: {},
    [ALARM_FIELDS.pbehaviorInfoReason]: {},
    [ALARM_FIELDS.pbehaviorInfoType]: {},
    [ALARM_FIELDS.pbehaviorInfoCanonicalType]: {},
  };

  return ADVANCED_SEARCH_OPTIONS_MAP;
};

export const useAdvancedSearchFieldVariables = () => {
  const { t, tc } = useI18n();

  const options = useAlarmAdvancedSearchVariables();

  return computed(() => Object.entries(ALARM_GROUPED_ADVANCED_SEARCH_FIELDS).reduce((acc, [group, items]) => {
    const header = t(`advancedSearch.groups.${group}`);

    acc.push(
      { header, value: header },
      ...items.map(field => ({
        ...options[field],
        value: field,
        text: tc(ALARM_FIELDS_TO_LABELS_KEYS[field], 2),
      })),
    );

    return acc;
  }, []));
};

const groupedAdvancedSearchFieldsToArray = () => (
  Object.entries(ALARM_GROUPED_ADVANCED_SEARCH_FIELDS).reduce((acc, [group, items]) => {
    acc.push({ header: group }, ...items);

    return acc;
  }, {})
);

const GROUP_TYPES = {
  field: 'field',
  union: 'union',
};

const advancedSearchGroupToForm = (group = {}) => ({
  key: uid(),
  field: group.field,
  operator: group.operator,
  valueType: group.valueType,
  value: group.value,
  secondValue: group.secondValue,
});

export default {
  components: { AdvancedSearchGroup, AdvancedSearchChip, AdvancedSearchMultipleChip },
  mixins: [Themeable],
  props: {
    isMenuActive: {
      type: Boolean,
      default: false,
    },
    menuProps: {
      type: Object,
      default: () => ({
        // openOnClick: false,
        disableKeys: true,
        closeOnContentClick: false,
        ignoreClickOutsideOnActivator: true,
        maxHeight: 304,
        nudgeBottom: 1,
        bottom: true,
        offsetY: true,
        transition: false,
      }),
    },
  },
  setup(props) {
    const inputValue = ref('');

    const groups = ref([advancedSearchGroupToForm()]);
    const activeGroup = ref(0);

    const localValue = ref([]);
    const step = ref(ADVANCED_SEARCH_STEP_TYPES.field);
    const activeKey = ref(null);

    const wrapperClass = computed(() => ({
      'v-input--is-focused': props.isMenuActive,
    }));

    const items = computed(() => groupedAdvancedSearchFieldsToArray());

    const fields = useAdvancedSearchFieldVariables();

    // const UNTYPICAL_STEP_TYPES_BY_FIELDS_AND_OPERATORS = { ...generateUntypicalDateTypes() };

    const nextStep = () => {
      switch (step.value) {
        case ADVANCED_SEARCH_STEP_TYPES.field:
          step.value = ADVANCED_SEARCH_STEP_TYPES.operator;
          return;
        case ADVANCED_SEARCH_STEP_TYPES.operator:
          step.value = ADVANCED_SEARCH_STEP_TYPES.value;
          // TODO: may be: ADVANCED_SEARCH_STEPS.valueType, ADVANCED_SEARCH_STEPS.value, ADVANCED_SEARCH_STEPS.union
          break;
        case ADVANCED_SEARCH_STEP_TYPES.valueType:
          step.value = ADVANCED_SEARCH_STEP_TYPES.value;
          return;
        case ADVANCED_SEARCH_STEP_TYPES.value:
          step.value = ADVANCED_SEARCH_STEP_TYPES.union;
          // TODO: may be: ADVANCED_SEARCH_STEPS.union, ADVANCED_SEARCH_STEPS.valueSecond
          break;
        case ADVANCED_SEARCH_STEP_TYPES.valueSecond:
          step.value = ADVANCED_SEARCH_STEP_TYPES.union;
          return;
        case ADVANCED_SEARCH_STEP_TYPES.union:
          step.value = ADVANCED_SEARCH_STEP_TYPES.field;
      }
    };

    const prev = () => {
      step.value = ADVANCED_SEARCH_STEP_TYPES.field;

      const fieldIndex = localValue.value.findLastIndex(({ type }) => type === ADVANCED_SEARCH_STEP_TYPES.field);

      localValue.value.splice(fieldIndex || 0);
    };

    const isUnion = computed(() => activeGroup.value % 2 !== 0);

    const selectItem = (item) => {
      console.log(item);
    };

    const addItem = (item) => {
      localValue.value.push({ key: uid(), value: item, type: step.value });

      nextStep();
    };

    const makeActive = (item) => {
      activeKey.value = item.key;
    };

    const resetActive = (item) => {
      console.log(item.key, activeKey.value);
      if (item.key !== activeKey.value) {
        return;
      }

      activeKey.value = null;
    };

    const update = (value, index) => {
      set(groups.value, index, value);
    };

    return {
      groups,
      fields,
      inputValue,
      activeKey,
      localValue,
      step,
      items,
      wrapperClass,

      addItem,
      selectItem,
      makeActive,
      resetActive,
      update,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-new-advanced-search { // TODO: remove new
  --v-chip-gap: 4px;
  --input-min-inline-size: 20ch;

  &::v-deep {
    input {
      flex: 0 1 auto;
      field-sizing: content;
      min-inline-size: var(--input-min-inline-size);
    }

    .layout {
      padding: var(--v-chip-gap);
      gap: var(--v-chip-gap);
    }

    .v-chip {
      padding: 0 8px;
      margin: 0;

      &:has(> .v-chip__content > .v-chip) {
        padding: 0 6px !important;

        button {
          margin: 0 -2px 0 0 !important;
        }
      }

      &__content {
        gap: var(--v-chip-gap);
      }

      .v-chip {
        height: 24px;

        &.theme--light {
          background: var(--v-application-background-base, #FFFFFF);
        }

        &.theme--dark {
          background: var(--v-application-background-base, #121212);
        }
      }
    }

    button {
      margin-left: 4px !important;
    }
  }
}
</style>
