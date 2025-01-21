<template>
  <div
    :class="[themeClasses]"
    class="c-new-advanced-search v-input v-input--hide-details theme--light
    v-text-field v-text-field--single-line v-text-field--is-booted v-select v-autocomplete primary--text"
  >
    <div class="v-input__control">
      <div class="v-input__slot">
        <div class="v-text-field__slot">
          <v-layout class="c-new-advanced-search__groups-wrapper gap-1" align-center wrap>
            <advanced-search-group
              v-for="(rule, index) in rules"
              v-model="rules[index]"
              :key="rule.key"
              :attributes="items"
              :active-key="activeKey"
              :union="index % 2 === 1"
              @input="update($event, index)"
              @click:chip="makeActive"
              @focusout="resetActive"
              @next="next($event, index)"
              @remove="remove(index)"
            />
          </v-layout>
        </div>
        <div class="v-input__append-inner">
          <v-menu bottom>
            <template #activator="{ on }">
              <c-action-btn
                :tooltip="$t('common.search')"
                icon="history"
                v-on="on"
              />
            </template>
          </v-menu>
        </div>
      </div>
    </div>
    <div class="v-input__append-outer">
      <div class="v-input__icon v-input__icon--append-outer">
        <c-action-btn
          :tooltip="$t('common.search')"
          icon="search"
          @click="submit"
        />
        <c-action-btn
          :tooltip="$t('common.clearSearch')"
          icon="clear"
          @click="clear"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { computed, ref, set } from 'vue';
import Themeable from 'vuetify/lib/mixins/themeable';

import { ALARM_FIELDS, PATTERN_OPERATORS, PATTERN_QUICK_RANGES, PATTERN_STRING_OPERATORS } from '@/constants';

import { advancedSearchRuleToForm } from '@/helpers/search/new-advanced-search';

import AdvancedSearchGroup from '@/components/common/search/partials/new/advanced-search-group.vue';
import { useComponentInstance } from '@/hooks/vue';
import { useEntity } from '@/hooks/store/modules/entity';

const ALARM_ENTITY_FIELDS_PREFIX = 'entity';
const ALARM_PBEHAVIOR_FIELDS_PREFIX = 'v.pbehavior_info';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: { AdvancedSearchGroup },
  mixins: [Themeable],
  props: {
    searches: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const { fetchContextEntitiesListWithoutStore } = useEntity();
    const instance = useComponentInstance();
    const rules = ref([advancedSearchRuleToForm()]);
    const activeKey = ref();

    const hasAlarmField = computed(() => rules.value.some(({ attribute }) => !attribute.startsWith(ALARM_PBEHAVIOR_FIELDS_PREFIX) && !attribute.startsWith(ALARM_PBEHAVIOR_FIELDS_PREFIX)));
    const hasEntityField = computed(() => rules.value.some(({ attribute }) => attribute.startsWith(ALARM_ENTITY_FIELDS_PREFIX)));
    const hasPbehaviorField = computed(() => rules.value.some(({ attribute }) => attribute.startsWith(ALARM_PBEHAVIOR_FIELDS_PREFIX)));

    const items = computed(() => [
      { header: 'Common' },
      {
        value: ALARM_FIELDS.displayName,
        operators: [
          ...PATTERN_STRING_OPERATORS,

          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
        ],
        text: 'Display name',
      },
      {
        value: ALARM_FIELDS.output,
        operators: [
          ...PATTERN_STRING_OPERATORS,

          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
        ],
        text: 'Output',
      },
      {
        value: ALARM_FIELDS.creationDate,
        text: 'Creation date',
        operators: PATTERN_QUICK_RANGES,
      },
      {
        value: ALARM_FIELDS.component,
        operators: [
          ...PATTERN_STRING_OPERATORS,

          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
        ],
        fetchValues: fetchContextEntitiesListWithoutStore,
        text: 'Component',
      },
    ]);

    const update = (value, index) => {
      instance.errors.clear();
      set(rules.value, index, value);
    };

    const makeActive = (key) => {
      setTimeout(() => activeKey.value = key); // TODO: refactor it
    };

    const resetActive = () => {
      activeKey.value = null;
    };

    const next = (nextStep, index) => {
      console.log(nextStep, rules.value[index + 1]);
      if (!nextStep && !rules.value[index + 1]) {
        const newRule = advancedSearchRuleToForm();

        rules.value.push(newRule);
        activeKey.value = `${newRule.key}.attribute`;
        return;
      }

      activeKey.value = !nextStep ? `${rules.value[index + 1].key}.attribute` : `${rules.value[index].key}.${nextStep}`;
    };

    const remove = (index) => {
      if (index === rules.value.length - 1) {
        set(rules.value, index, advancedSearchRuleToForm());

        return;
      }

      rules.value.splice(index, index === rules.value.length - 2 ? 1 : 2);
    };

    const submit = async () => {
      const isValid = await instance.$validator.validateAll();

      if (isValid) {
        emit('submit', rules.value);
      }
    };

    const clear = () => {
      instance.errors.clear();
      rules.value = [advancedSearchRuleToForm()];
    };

    return {
      rules,
      activeKey,

      items,

      update,
      makeActive,
      resetActive,
      next,
      remove,
      submit,
      clear,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-new-advanced-search { // TODO: remove new
  --v-chip-gap: 4px;
  --input-min-inline-size: 20ch;

  &__groups-wrapper > * {
    flex: 0 1 auto;
  }

  &::v-deep {
    input {
      flex: 0 1 auto;
      field-sizing: content;
      min-inline-size: var(--input-min-inline-size);
    }

    .layout {
      padding: var(--v-chip-gap) 0;
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
