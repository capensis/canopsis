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
  },
  setup(props) {
    const { t, tc } = useI18n();

    const availableActions = computed(() => {
      const { actions, featuresActionsKey } = {
        [props.massive]: {
          actions: UNIQUE_ALARM_LIST_ACTIONS_TYPES_TO_LABELS_KEYS,
          featuresActionsKey: 'components.alarmListActionPanel.computed.actions',
        },
        [!props.massive]: {
          actions: UNIQUE_ALARM_LIST_MASS_ACTIONS_TYPES_TO_LABELS_KEYS,
          featuresActionsKey: 'components.alarmListMassActionsPanel.computed.actions',
        },
      }.true;

      const preparedActions = Object.entries(actions).map(([value, textKey]) => ({
        value,
        text: t(textKey),
      }));

      let featuresActions = [];

      try {
        const contextForFeatureSwitcher = {
          checkAccess: () => true,
          $t: t,
          $tc: tc,
        };

        featuresActions = featuresService.has(featuresActionsKey)
        && featuresService.call(featuresActionsKey, contextForFeatureSwitcher, []);

        if (featuresActions?.length) {
          const preparedFeaturesActions = featuresActions.map(({ type, title }) => ({ value: type, text: title }));

          preparedActions.unshift(...preparedFeaturesActions);
        }
      } catch (err) {
        console.error(err);
      }

      return preparedActions;
    });

    return {
      availableActions,
    };
  },
};
</script>
