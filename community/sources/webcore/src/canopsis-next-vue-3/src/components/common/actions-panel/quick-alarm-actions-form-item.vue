<template>
  <c-card-iterator-item
    :drag-handle-class="dragHandleClass"
    small
    @remove="$emit('remove')"
  >
    <template #header>
      <c-select
        :model-value="modelValue"
        :items="availableActions"
        :label="t('common.action', 1)"
        :name="name"
        item-title="text"
        item-value="value"
        hide-details
        required
        @update:model-value="$emit('update:modelValue', $event)"
      />
    </template>
  </c-card-iterator-item>
</template>

<script>
import { computed } from 'vue';

import {
  UNIQUE_ALARM_LIST_ACTIONS_TYPES_TO_LABELS_KEYS,
  UNIQUE_ALARM_LIST_MASS_ACTIONS_TYPES_TO_LABELS_KEYS,
} from '@/constants';

import { featuresService } from '@/services/features';

import { useI18n } from '@/hooks/i18n';

export default {
  name: 'QuickAlarmActionsFormItem',

  props: {
    modelValue: {
      type: String,
      required: true,
    },
    dragHandleClass: {
      type: String,
      required: true,
    },
    massive: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'action',
    },
    selectedActions: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue', 'remove'],

  setup(props) {
    const { t } = useI18n();

    const availableActions = computed(() => {
      const { actions, featuresActionsKey } = props.massive
        ? {
            actions: UNIQUE_ALARM_LIST_MASS_ACTIONS_TYPES_TO_LABELS_KEYS,
            featuresActionsKey: 'components.alarmListMassActionsPanel.computed.actions',
          }
        : {
            actions: UNIQUE_ALARM_LIST_ACTIONS_TYPES_TO_LABELS_KEYS,
            featuresActionsKey: 'components.alarmListActionPanel.computed.actions',
          };

      const preparedActions = Object.entries(actions).map(([value, textKey]) => ({
        value,
        text: t(textKey),
      }));

      try {
        const contextForFeatureSwitcher = {
          checkAccess: () => true,
          $t: t,
          $tc: t,
        };

        const featuresActions = featuresService.has(featuresActionsKey)
          ? featuresService.call(featuresActionsKey, contextForFeatureSwitcher, [])
          : [];

        if (featuresActions?.length) {
          preparedActions.unshift(
            ...featuresActions.map(({ type, title }) => ({ value: type, text: title })),
          );
        }
      } catch (err) {
        console.error(err);
      }

      return preparedActions.filter(
        preparedAction =>
          preparedAction.value === props.modelValue ||
          !props.selectedActions.some(action => action.value === preparedAction.value),
      );
    });

    return {
      t,
      availableActions,
    };
  },
};
</script>
