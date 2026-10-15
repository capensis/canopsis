<template>
  <v-select
    v-validate="'required'"
    :value="value"
    :items="preparedTriggers"
    :disabled="disabled"
    :label="label || $tc('common.trigger', 2)"
    :error-messages="errorMessages"
    :name="name"
    item-value="type"
    item-disabled="deprecated"
    multiple
    chips
    return-object
    @change="changeValue"
  >
    <template #selection="{ item, index }">
      <v-tooltip
        :disabled="!item.deprecated || errors.has(getChipName(index))"
        top
      >
        <template #activator="{ on }">
          <v-chip
            v-validate=""
            :class="getSelectedClass(item)"
            :close="errors.has(getChipName(index)) || item.deprecated"
            :name="`${name}.${index}.type`"
            :color="errors.has(getChipName(index)) ? 'error' : ''"
            v-on="on"
            @click:close="removeItemFromArray(index)"
          >
            {{ getSelectedText(item) }}
          </v-chip>
        </template>
        <span>
          {{
            errors.has(getChipName(index))
              ? errors.collect(getChipName(index))
              : item.deprecated
                ? $t('common.deprecatedTrigger')
                : ''
          }}
        </span>
      </v-tooltip>
    </template>
    <template #item="{ item, attrs, on, parent }">
      <v-list-item
        v-bind="attrs"
        :active-class="errors.has(getAdditionalValueFieldName(item.type)) ? 'error--text' : attrs.activeClass"
        @click="on.click"
      >
        <v-list-item-action>
          <v-checkbox
            :input-value="attrs.inputValue"
            :color="parent.color"
            hide-details
          />
        </v-list-item-action>
        <v-list-item-content>
          <v-layout
            class="fill-width"
            align-center
            justify-space-between
          >
            <v-flex>{{ item.text }}</v-flex>
            <component
              v-if="additionalValuesComponentsByTypes[item.type]"
              v-bind="additionalValuesComponentsByTypes[item.type].bind"
              :is="additionalValuesComponentsByTypes[item.type].is"
              :disabled="!attrs.inputValue"
              class="ml-1"
              v-on="additionalValuesComponentsByTypes[item.type].on"
              @click.prevent.stop=""
            />
          </v-layout>
        </v-list-item-content>
        <v-list-item-action v-if="item.helpText">
          <c-help-icon
            :text="item.helpText"
            color="info"
            size="20"
            top
          />
        </v-list-item-action>
      </v-list-item>
    </template>
  </v-select>
</template>

<script>
import {
  computed,
  ref,
  watch,
  set,
  nextTick,
  onBeforeUnmount,
} from 'vue';
import { find } from 'lodash';

import { TRIGGERS_TYPES, PRO_TRIGGERS } from '@/constants';

import { setSeveralFields } from '@/helpers/immutable';
import { isDeprecatedTrigger } from '@/helpers/entities/scenario/form';

import { useValidator } from '@/hooks/validator/validator';
import { useI18n } from '@/hooks/i18n';
import { useInfo } from '@/hooks/store/modules/info';
import { useComponentInstance } from '@/hooks/vue';
import { useArrayModelField } from '@/hooks/form/array-model-field';

const DEFAULT_TRANSLATION_PREFIX = 'common.triggers';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'triggers',
    },
    withAdditionalValues: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    types: {
      type: Array,
      default: () => Object.values(TRIGGERS_TYPES),
    },
    translationKeyPrefix: {
      type: String,
      default: DEFAULT_TRANSLATION_PREFIX,
    },
  },
  setup(props, { emit }) {
    const additionalValuesKeysByTriggers = {
      [TRIGGERS_TYPES.eventscount]: 'threshold',
    };

    const validator = useValidator();
    const { t, tc, te } = useI18n();
    const { updateModel, removeItemFromArray } = useArrayModelField(props, emit);
    const { isProVersion, eventsCountTriggerDefaultThreshold } = useInfo();
    const vm = useComponentInstance();

    const preparedTriggersByTypes = ref({});

    const deprecatedValues = computed(() => props.value.filter(({ type }) => isDeprecatedTrigger(type)));

    const defaultAdditionalValuesByTriggers = computed(() => ({
      [TRIGGERS_TYPES.eventscount]: eventsCountTriggerDefaultThreshold.value ?? '',
    }));

    const preparedTriggers = computed(() => props.types.reduce((acc, type) => {
      if (!PRO_TRIGGERS.includes(type) || isProVersion.value) {
        const key = `${props.translationKeyPrefix}.${type}`;
        const { text, helpText } = props.translationKeyPrefix !== DEFAULT_TRANSLATION_PREFIX && te(key)
          ? t(key)
          : t(`${DEFAULT_TRANSLATION_PREFIX}.${type}`);

        acc.push({
          ...(preparedTriggersByTypes.value[type] ?? {}),

          text,
          helpText,
          deprecated: isDeprecatedTrigger(type),
        });
      }

      return acc;
    }, []));

    const errorMessages = computed(() => validator.errors.collect(props.name, null, false)
      .map((item) => {
        const messageMap = {
          max_value: tc(
            'errors.triggerMustNotUsed',
            deprecatedValues.value.length,
            { field: deprecatedValues.value.join(', ') },
          ),
        };

        return messageMap[item.rule] ?? item.msg;
      }));

    /**
     * Get name property for chip field
     *
     * @param {number} index
     * @return {string}
     */
    const getChipName = index => `${props.name}.${index}.type`;

    /**
     * Get name property for additional value field
     *
     * @param {string} type
     * @return {string}
     */
    const getAdditionalValueFieldName = type => `${props.name}.${type}.additionalValue`;

    /**
     * Get CSS class for selected item
     *
     * @param {boolean} deprecated
     * @param {string} type
     * @return {{ error: boolean, 'error--text': boolean }}
     */
    const getSelectedClass = ({ deprecated, type }) => ({
      'error--text': deprecated,
      error: validator.errors.has(getAdditionalValueFieldName(type)),
    });

    /**
     * Get text for selected item
     *
     * @param {string} type
     * @param {string} text
     * @param {number | string | boolean} [additionalValue]
     * @return {VueI18n.TranslateResult|*}
     */
    const getSelectedText = ({ type, text, [additionalValuesKeysByTriggers[type]]: additionalValue } = {}) => {
      const messageKey = `${props.translationKeyPrefix}.${type}.selectedText`;

      return te(messageKey)
        ? t(messageKey, { additionalValue })
        : text;
    };

    /**
     * Set preparedTriggersByTypes by types array
     *
     * @param {string[]} [types = []]
     */
    const setPreparedTriggersByTypes = (types = []) => {
      preparedTriggersByTypes.value = types.reduce((acc, type) => {
        const additionalValueKey = additionalValuesKeysByTriggers[type];
        const {
          [additionalValueKey]: additionalValue = defaultAdditionalValuesByTriggers.value[type],
        } = find(props.value, { type }) ?? {};

        const newTypeValue = { type };

        if (additionalValueKey) {
          newTypeValue[additionalValueKey] = additionalValue;
        }

        acc[type] = newTypeValue;

        return acc;
      }, {});
    };

    /**
     * Change value on select field handler
     *
     * @param {Trigger[]} value
     */
    const changeValue = (value = []) => {
      updateModel(value.map(({ type }) => preparedTriggersByTypes.value[type]));

      nextTick(() => validator.validate(props.name));
    };

    /**
     * Change additional value handler
     *
     * @param {string} type
     * @param {number | string | boolean} additionalValue
     */
    const changeAdditionalValue = (type, additionalValue) => {
      set(
        preparedTriggersByTypes.value[type],
        additionalValuesKeysByTriggers[type],
        additionalValue,
      );

      updateModel(props.value.map(trigger => (
        trigger.type === type
          ? preparedTriggersByTypes.value[type]
          : trigger
      )));
    };

    /**
     * Attach rule for deprecatedValues checking into validator
     */
    const attachMaxValueRule = () => {
      validator.attach({
        name: props.name,
        rules: 'max_value:0',
        getter: () => deprecatedValues.value.length,
        vm,
      });
    };

    /**
     * Detach rule for deprecatedValues checking from validator
     */
    const detachRules = () => {
      validator.detach(props.name);
    };

    const additionalValuesComponentsByTypes = computed(() => {
      if (!props.withAdditionalValues) {
        return {};
      }

      const additionalValuesComponents = [
        {
          type: TRIGGERS_TYPES.eventscount,
          is: 'c-number-field',
          bind: {
            class: ['mt-0', 'pt-2'],
            hideDetails: true,
            min: 1,
            required: true,
          },
        },
      ];

      return additionalValuesComponents.reduce((acc, { type, ...rest }) => {
        const trigger = preparedTriggersByTypes.value[type];
        const name = getAdditionalValueFieldName(type);

        if (!trigger) {
          return acc;
        }

        acc[type] = setSeveralFields(rest, {
          'bind.name': name,
          'bind.value': trigger[additionalValuesKeysByTriggers[type]],
          'bind.label': t(`common.triggers.${type}.additionalFieldLabel`),
          'bind.errorMessages': validator.errors.collect(name),
          'on.input': () => value => changeAdditionalValue(type, value),
        });

        return acc;
      }, {});
    });

    watch(() => props.types, setPreparedTriggersByTypes, { immediate: true });

    attachMaxValueRule();

    onBeforeUnmount(() => {
      detachRules();
    });

    return {
      preparedTriggers,
      errorMessages,
      additionalValuesComponentsByTypes,
      getChipName,
      getAdditionalValueFieldName,
      getSelectedClass,
      getSelectedText,
      changeValue,
      removeItemFromArray,
    };
  },
};
</script>
