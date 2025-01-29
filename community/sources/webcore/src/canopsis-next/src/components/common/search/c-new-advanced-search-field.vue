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
              :first="index === 0"
              :allow-or="allowOr"
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

import { advancedSearchRulesToForm, formToAdvancedSearchRules, advancedSearchItemToForm } from '@/helpers/search/new-advanced-search';

import { useComponentInstance } from '@/hooks/vue';
import { useEntity } from '@/hooks/store/modules/entity';

import AdvancedSearchGroup from '@/components/common/search/partials/new/advanced-search-group.vue';

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
    const rules = ref([advancedSearchItemToForm()]);
    const activeKey = ref();

    const hasEntityField = computed(() => (
      rules.value.some(({ attribute }) => attribute.startsWith(ALARM_ENTITY_FIELDS_PREFIX))
    ));
    const hasPbehaviorField = computed(() => (
      rules.value.some(({ attribute }) => attribute.startsWith(ALARM_PBEHAVIOR_FIELDS_PREFIX))
    ));
    const hasAlarmField = computed(() => (
      rules.value.some(({ attribute }) => (
        !attribute.startsWith(ALARM_PBEHAVIOR_FIELDS_PREFIX)
        && !attribute.startsWith(ALARM_PBEHAVIOR_FIELDS_PREFIX)
      ))
    ));
    const hasOr = computed(() => rules.value.some(({ union }) => union === 'OR'));

    const allowOr = computed(() => [
      hasEntityField.value,
      hasPbehaviorField.value,
      hasAlarmField.value,
    ].filter(Boolean).length <= 1);
    const allowAlarmFields = computed(() => !hasOr.value || (!hasEntityField.value && !hasPbehaviorField.value));
    const allowEntityFields = computed(() => !hasOr.value || (!hasAlarmField.value && !hasPbehaviorField.value));
    const allowPbehaviorFields = computed(() => !hasOr.value || (!hasAlarmField.value && !hasEntityField.value));

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
        disabled: !allowAlarmFields.value,
      },
      {
        value: ALARM_FIELDS.entityInfos,
        text: 'Entity infos',
        disabled: !allowEntityFields.value,
        items: [
          {
            value: 'dictionary1',
            text: 'dictionary1',
            items: [{ value: 'value', text: 'Value' }, { value: 'name', text: 'Name' }],
          },
          {
            value: 'dictionary2',
            text: 'dictionary2',
            items: [{ value: 'value', text: 'Value' }, { value: 'name', text: 'Name' }],
          },
        ],
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
      if (!nextStep && !rules.value[index + 1]) {
        const newRule = advancedSearchItemToForm();

        console.log(newRule, `${newRule.key}.${index % 2 === 0 ? 'union' : 'attribute'}`);

        rules.value.push(newRule);
        activeKey.value = `${newRule.key}.${index % 2 === 0 ? 'union' : 'attribute'}`;
        return;
      }

      activeKey.value = !nextStep ? `${rules.value[index + 1].key}.attribute` : `${rules.value[index].key}.${nextStep}`;
    };

    const remove = (index) => {
      if (index === rules.value.length - 1) {
        set(rules.value, index, advancedSearchItemToForm());

        return;
      }

      rules.value.splice(index, index === rules.value.length - 2 ? 1 : 2);
    };

    const submit = async () => {
      const isValid = await instance.$validator.validateAll();
      const json = '[{"key":"0e79fcb6573","attribute":"v.display_name","operator":"equal","field":"","fieldType":"string","dictionary":"","value":"fdsfdsf","range":"last1Hour","duration":{"value":1,"unit":"s"},"filled":["attribute","operator","value"],"rangeValue":{"from":0,"to":0},"union":null},{"key":"e79fcb65731","attribute":"","operator":"","field":"","fieldType":"string","dictionary":"","value":"","range":"last1Hour","duration":{"value":1,"unit":"s"},"filled":["union"],"rangeValue":{"from":0,"to":0},"union":"OR"},{"key":"79fcb657318","attribute":"v.output","operator":"contains","field":"","fieldType":"string","dictionary":"","value":"gfdgdfg","range":"last1Hour","duration":{"value":1,"unit":"s"},"filled":["attribute","operator","value"],"rangeValue":{"from":0,"to":0},"union":null},{"key":"9fcb6573188","attribute":"","operator":"","field":"","fieldType":"string","dictionary":"","value":"","range":"last1Hour","duration":{"value":1,"unit":"s"},"filled":[],"rangeValue":{"from":0,"to":0},"union":null}]';
      const json2 = '{"positions":["alarm_pattern","alarm_pattern"],"alarm_pattern":[[{"field":"v.display_name","cond":{"value":"fdsfdsf","type":"eq"}}],[{"field":"v.output","cond":{"value":"gfdgdfg","type":"contain"}}]],"entity_pattern":[],"pbehavior_pattern":[]}';
      console.log(advancedSearchRulesToForm(JSON.parse(json2)));
      // console.log(advancedSearchRulesToForm(JSON.parse(json2)));

      if (isValid) {
        emit('submit', formToAdvancedSearchRules(rules.value));
      }
    };

    const clear = () => {
      instance.errors.clear();
      rules.value = [advancedSearchRulesToForm()];
    };

    return {
      allowOr,
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

      // TODO: remove code below
      hasOr,
      hasEntityField,
      hasPbehaviorField,
      hasAlarmField,

      allowAlarmFields,
      allowEntityFields,
      allowPbehaviorFields,
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
