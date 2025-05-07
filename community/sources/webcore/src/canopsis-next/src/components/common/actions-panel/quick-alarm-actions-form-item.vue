<template>
  <c-card-iterator-item
    :drag-handle-class="dragHandleClass"
    small
    @remove="$emit('remove')"
  >
    <template #header>
      <c-select-field
        v-field="value"
        :items="availableActions"
        :label="$tc('common.action', 1)"
        :name="name"
        hide-details
        required
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
    value: {
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
  setup(props) {
    const { t, tc } = useI18n();

    const availableActions = computed(() => {
      // Determine which set of actions to use based on whether it's a massive action or not
      const { actions, featuresActionsKey } = {
        [!props.massive]: {
          // Use unique alarm list actions for non-massive operations
          actions: UNIQUE_ALARM_LIST_ACTIONS_TYPES_TO_LABELS_KEYS,
          featuresActionsKey: 'components.alarmListActionPanel.computed.actions',
        },
        [props.massive]: {
          // Use mass actions for massive operations
          actions: UNIQUE_ALARM_LIST_MASS_ACTIONS_TYPES_TO_LABELS_KEYS,
          featuresActionsKey: 'components.alarmListMassActionsPanel.computed.actions',
        },
      }.true;

      // Transform actions into a format suitable for display (value-text pairs)
      const preparedActions = Object.entries(actions).map(([value, textKey]) => ({
        value,
        text: t(textKey),
      }));

      let featuresActions = [];

      try {
        // Create context for feature switcher with necessary methods
        const contextForFeatureSwitcher = {
          checkAccess: () => true,
          $t: t,
          $tc: tc,
        };

        // Check if feature exists and get its actions
        featuresActions = featuresService.has(featuresActionsKey)
        && featuresService.call(featuresActionsKey, contextForFeatureSwitcher, []);

        // If feature actions exist, add them to the beginning of the actions list
        if (featuresActions?.length) {
          const preparedFeaturesActions = featuresActions.map(({ type, title }) => ({ value: type, text: title }));

          preparedActions.unshift(...preparedFeaturesActions);
        }
      } catch (err) {
        console.error(err);
      }

      // Filter out actions that are already selected, except for the current value
      return preparedActions.filter(preparedAction => (
        preparedAction.value === props.value || !props.selectedActions.includes(preparedAction.value)
      ));
    });

    return {
      availableActions,
    };
  },
};
</script>
