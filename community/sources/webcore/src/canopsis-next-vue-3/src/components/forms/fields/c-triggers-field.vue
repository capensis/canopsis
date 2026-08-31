<template>
  <c-select
    :model-value="selectedTriggers"
    :items="preparedTriggers"
    :disabled="disabled"
    :label="label || $t('common.trigger', 2)"
    :name="name"
    :error-messages="errorMessages"
    item-value="type"
    item-title="text"
    multiple
    chips
    return-object
    standalone
    @update:model-value="changeValue"
  >
    <template #selection="{ item, index }">
      <v-tooltip
        :disabled="!getPreparedTrigger(item).deprecated || hasError(getChipName(index))"
        location="top"
      >
        <template #activator="{ props: tooltipProps }">
          <v-chip
            v-bind="tooltipProps"
            :class="getSelectedClass(getPreparedTrigger(item))"
            :closable="hasError(getChipName(index)) || getPreparedTrigger(item).deprecated"
            :name="`${name}.${index}.type`"
            :color="hasError(getChipName(index)) ? 'error' : ''"
            @click:close="removeItemFromArray(index)"
          >
            {{ getSelectedText(getPreparedTrigger(item)) }}
          </v-chip>
        </template>
        <span>
          {{
            hasError(getChipName(index))
              ? collectErrors(getChipName(index))
              : getPreparedTrigger(item).deprecated
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
          hasError(getAdditionalValueFieldName(getPreparedTrigger(item).type))
            ? 'error--text'
            : itemProps.activeClass
        "
      >
        <template #prepend>
          <v-checkbox-btn
            :model-value="isTriggerSelected(getPreparedTrigger(item))"
            color="primary"
            hide-details
          />
        </template>

        <v-list-item-title>
          <c-layout align-center justify-space-between>
            <span>{{ getPreparedTrigger(item).text }}</span>
            <component
              v-if="additionalValuesComponentsByTypes[getPreparedTrigger(item).type]"
              v-bind="additionalValuesComponentsByTypes[getPreparedTrigger(item).type].bind"
              :is="additionalValuesComponentsByTypes[getPreparedTrigger(item).type].is"
              :disabled="!isTriggerSelected(getPreparedTrigger(item))"
              class="ml-1"
              v-on="additionalValuesComponentsByTypes[getPreparedTrigger(item).type].on"
              @click.prevent.stop=""
            />
          </c-layout>
        </v-list-item-title>

        <template v-if="getPreparedTrigger(item).helpText" #append>
          <c-help-icon
            :text="getPreparedTrigger(item).helpText"
            color="info"
            size="20"
            location="top"
          />
        </template>
      </v-list-item>
    </template>
  </c-select>
</template>

<script>
import { find } from 'lodash';
import { computed, ref, watch } from 'vue';

import { TRIGGERS_TYPES, PRO_TRIGGERS } from '@/constants';

import { isDeprecatedTrigger } from '@/helpers/entities/scenario/form';
import { setSeveralFields } from '@/helpers/immutable';

import { useField, useFormErrorsMap } from '@/hooks/form';
import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useI18n } from '@/hooks/i18n';
import { useInfo } from '@/hooks/store/modules/info';

const DEFAULT_TRANSLATION_PREFIX = 'common.triggers';

export default {
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

    const deprecatedValues = computed(() => props.modelValue.filter(({ type }) => isDeprecatedTrigger(type)),
    );

    const validateTriggers = (value = []) => {
      if (!value.length) {
        return t('validation.messages.required');
      }

      if (deprecatedValues.value.length) {
        return t('errors.triggerMustNotUsed', deprecatedValues.value.length, {
          field: deprecatedValues.value.map(({ type }) => type).join(', '),
        });
      }

      return true;
    };
    const { setValue: setTriggersFieldValue } = useField(
      computed(() => props.name),
      validateTriggers,
      {
        initialValue: props.modelValue,
        validateOnValueUpdate: true,
      },
    );
    const { collectErrors, hasError } = useFormErrorsMap();

    const defaultAdditionalValuesByTriggers = computed(() => ({
      [TRIGGERS_TYPES.eventscount]: eventsCountTriggerDefaultThreshold.value ?? '',
    }));

    const getTriggerTranslationKey = (type) => {
      const key = `${props.translationKeyPrefix}.${type}`;

      return props.translationKeyPrefix !== DEFAULT_TRANSLATION_PREFIX && te(`${key}.text`)
        ? key
        : `${DEFAULT_TRANSLATION_PREFIX}.${type}`;
    };

    const triggerToPreparedTrigger = (trigger = {}) => {
      const { type } = trigger;
      const triggerKey = getTriggerTranslationKey(type);

      return {
        ...trigger,

        text: te(`${triggerKey}.text`) ? t(`${triggerKey}.text`) : type,
        helpText: te(`${triggerKey}.helpText`) ? t(`${triggerKey}.helpText`) : '',
        deprecated: isDeprecatedTrigger(type),
      };
    };

    const preparedTriggers = computed(() => props.types.reduce((acc, type) => {
        if (!PRO_TRIGGERS.includes(type) || isProVersion.value) {
          acc.push(triggerToPreparedTrigger({
            ...preparedTriggersByTypes.value[type],
            type,
          }));
        }

        return acc;
      }, []),
    );
    const selectedTriggers = computed(() => props.modelValue.map((trigger) => {
        const preparedTrigger = preparedTriggers.value.find(({ type }) => type === trigger.type);

        return preparedTrigger
          ? {
              ...preparedTrigger,
              ...trigger,
            }
          : triggerToPreparedTrigger(trigger);
      }),
    );

    const errorMessages = computed(() => {
      if (!hasError(props.name)) {
        return [];
      }

      return collectErrors(props.name);
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
    const isTriggerSelected = trigger => props.modelValue.some(({ type }) => type === trigger.type);

    /**
     * Get CSS class for selected item
     *
     * @param {boolean} deprecated
     * @param {string} type
     * @return {{ error: boolean, 'error--text': boolean }}
     */
    const getSelectedClass = ({ deprecated, type }) => ({
      'error--text': deprecated,
      error: hasError(getAdditionalValueFieldName(type)),
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
     * Get prepared trigger by select slot item.
     *
     * @param {Object} item
     * @return {Trigger}
     */
    const getPreparedTrigger = (item = {}) => {
      const type = item.type ?? item.value;

      return preparedTriggers.value.find(triggerItem => triggerItem.type === type)
        ?? triggerToPreparedTrigger(item);
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
      const nextValue = value.map(({ type }) => preparedTriggersByTypes.value[type] ?? { type });

      updateModel(nextValue);
      setTriggersFieldValue(nextValue);
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
        props.modelValue.map(trigger => trigger.type === type ? preparedTriggersByTypes.value[type] : trigger,
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
            'bind.errorMessages': collectErrors(name),
          }),
          on: {
            'update:value': value => changeAdditionalValue(type, value),
          },
        };

        return acc;
      }, {});
    });

    watch(() => props.types, setPreparedTriggersByTypes, { immediate: true });
    watch(
      () => props.modelValue,
      value => setTriggersFieldValue(value, false),
    );

    return {
      preparedTriggers,
      selectedTriggers,
      errorMessages,
      additionalValuesComponentsByTypes,
      collectErrors,
      hasError,
      getChipName,
      getAdditionalValueFieldName,
      isTriggerSelected,
      getPreparedTrigger,
      getSelectedClass,
      getSelectedText,
      changeValue,
      removeItemFromArray,
    };
  },
};
</script>
