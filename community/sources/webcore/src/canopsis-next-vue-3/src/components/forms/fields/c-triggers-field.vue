<template>
  <v-select
    v-validate="'required'"
    :model-value="modelValue"
    :items="preparedTriggers"
    :disabled="disabled"
    :label="label || $t('common.trigger', 2)"
    :error-messages="errorMessages"
    :name="name"
    item-value="type"
    item-title="text"
    multiple
    chips
    return-object
    @update:model-value="changeValue"
  >
    <template #selection="{ item, index }">
      <v-tooltip :disabled="!item.raw.deprecated || errors.has(getChipName(index))" location="top">
        <template #activator="{ props: tooltipProps }">
          <v-chip
            v-validate=""
            v-bind="tooltipProps"
            :class="getSelectedClass(item.raw)"
            :closable="errors.has(getChipName(index)) || item.raw.deprecated"
            :name="`${name}.${index}.type`"
            :color="errors.has(getChipName(index)) ? 'error' : ''"
            @click:close="removeItemFromArray(index)"
          >
            {{ getSelectedText(item.raw) }}
          </v-chip>
        </template>
        <span>
          {{
            errors.has(getChipName(index))
              ? errors.collect(getChipName(index))
              : item.raw.deprecated
                ? $t('common.deprecatedTrigger')
                : ''
          }}
        </span>
      </v-tooltip>
    </template>
    <template #item="{ item, props: itemProps }">
      <v-list-item
        v-bind="itemProps"
        :active-class="
          errors.has(getAdditionalValueFieldName(item.raw.type))
            ? 'error--text'
            : itemProps.activeClass
        "
      >
        <template #prepend>
          <v-checkbox-btn :model-value="isTriggerSelected(item.raw)" color="primary" hide-details />
        </template>

        <v-list-item-title>
          <c-layout align-center justify-space-between>
            <span>{{ item.raw.text }}</span>
            <component
              v-if="additionalValuesComponentsByTypes[item.raw.type]"
              v-bind="additionalValuesComponentsByTypes[item.raw.type].bind"
              :is="additionalValuesComponentsByTypes[item.raw.type].is"
              :disabled="!isTriggerSelected(item.raw)"
              class="ml-1"
              v-on="additionalValuesComponentsByTypes[item.raw.type].on"
              @click.prevent.stop=""
            />
          </c-layout>
        </v-list-item-title>

        <template v-if="item.raw.helpText" #append>
          <c-help-icon
            :text="item.raw.helpText"
            color="info"
            size="20"
            location="top"
          />
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script>
import { find } from 'lodash';
import { computed, ref, watch, nextTick, onBeforeUnmount } from 'vue';

import { TRIGGERS_TYPES, PRO_TRIGGERS } from '@/constants';

import { isDeprecatedTrigger } from '@/helpers/entities/scenario/form';
import { setSeveralFields } from '@/helpers/immutable';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useI18n } from '@/hooks/i18n';
import { useInfo } from '@/hooks/store/modules/info';
import { useValidationAttachRule } from '@/hooks/validator/validation-attach-rule';

const DEFAULT_TRANSLATION_PREFIX = 'common.triggers';

export default {
  inject: ['$validator'],
  props: {
    modelValue: {
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
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const additionalValuesKeysByTriggers = {
      [TRIGGERS_TYPES.eventscount]: 'threshold',
    };

    const { t, te } = useI18n();
    const { updateModel, removeItemFromArray } = useArrayModelField(props, emit);
    const { isProVersion, eventsCountTriggerDefaultThreshold } = useInfo();

    const preparedTriggersByTypes = ref({});

    const deprecatedValues = computed(() =>
      props.modelValue.filter(({ type }) => isDeprecatedTrigger(type)),
    );

    const {
      validator,
      attachRule: attachMaxValueRule,
      detachRule: detachRules,
      validateRule: validateMaxValueRule,
    } = useValidationAttachRule(props.name, 'max_value:0', () => deprecatedValues.value.length);

    const defaultAdditionalValuesByTriggers = computed(() => ({
      [TRIGGERS_TYPES.eventscount]: eventsCountTriggerDefaultThreshold.value ?? '',
    }));

    const preparedTriggers = computed(() =>
      props.types.reduce((acc, type) => {
        if (!PRO_TRIGGERS.includes(type) || isProVersion.value) {
          const key = `${props.translationKeyPrefix}.${type}`;
          const { text, helpText } =
            props.translationKeyPrefix !== DEFAULT_TRANSLATION_PREFIX && te(key)
              ? t(key)
              : t(`${DEFAULT_TRANSLATION_PREFIX}.${type}`);

          acc.push({
            ...preparedTriggersByTypes.value[type],

            text,
            helpText,
            deprecated: isDeprecatedTrigger(type),
          });
        }

        return acc;
      }, []),
    );

    const errorMessages = computed(() => {
      if (!validator.errors.has(props.name)) {
        return [];
      }

      return [
        t('errors.triggerMustNotUsed', deprecatedValues.value.length, {
          field: deprecatedValues.value.map(({ type }) => type).join(', '),
        }),
      ];
    });

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
     * Check if trigger is selected
     *
     * @param {Trigger} trigger
     * @return {boolean}
     */
    const isTriggerSelected = trigger =>
      props.modelValue.some(({ type }) => type === trigger.type);

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
    const getSelectedText = ({
      type,
      text,
      [additionalValuesKeysByTriggers[type]]: additionalValue,
    } = {}) => {
      const messageKey = `${props.translationKeyPrefix}.${type}.selectedText`;

      return te(messageKey) ? t(messageKey, { additionalValue }) : text;
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
        } = find(props.modelValue, { type }) ?? {};

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

      nextTick(validateMaxValueRule);
    };

    /**
     * Change additional value handler
     *
     * @param {string} type
     * @param {number | string | boolean} additionalValue
     */
    const changeAdditionalValue = (type, additionalValue) => {
      preparedTriggersByTypes.value[type][additionalValuesKeysByTriggers[type]] = additionalValue;

      updateModel(
        props.modelValue.map(trigger =>
          trigger.type === type ? preparedTriggersByTypes.value[type] : trigger,
        ),
      );
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

        acc[type] = {
          ...setSeveralFields(rest, {
            'bind.name': name,
            'bind.value': trigger[additionalValuesKeysByTriggers[type]],
            'bind.label': t(`common.triggers.${type}.additionalFieldLabel`),
            'bind.errorMessages': validator.errors.collect(name),
          }),
          on: {
            'update:value': value => changeAdditionalValue(type, value),
          },
        };

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
      isTriggerSelected,
      getSelectedClass,
      getSelectedText,
      changeValue,
      removeItemFromArray,
    };
  },
};
</script>
