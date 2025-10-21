<template>
  <v-btn
    :disabled="alarmFilteringPending"
    color="primary"
    @click="runAlarmFiltering"
  >
    <v-progress-circular
      v-if="alarmFilteringPending"
      class="mr-2"
      size="16"
      width="2"
      indeterminate
    />
    {{ label }}
  </v-btn>
</template>

<script>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import { useAlarmFiltering } from '@/components/other/pbehavior/pbehaviors/hooks/alarm-filtering';

export default {
  setup(_, { emit }) {
    const { t } = useI18n();
    const modals = useModals();

    const refresh = () => emit('refresh');

    const {
      pending: alarmFilteringPending,
      runAlarmFiltering: runAlarmFilteringAction,
    } = useAlarmFiltering(refresh);

    const label = computed(() => (
      alarmFilteringPending.value
        ? t('pbehavior.alarmFilteringInProgress')
        : t('pbehavior.runAlarmFiltering')
    ));

    const runAlarmFiltering = () => modals.show({
      name: MODALS.confirmation,
      config: {
        title: t('modals.confirmationRunAlarmFiltering.title'),
        text: t('modals.confirmationRunAlarmFiltering.text'),
        action: runAlarmFilteringAction,
      },
    });

    return {
      runAlarmFiltering,
      alarmFilteringPending,
      label,
    };
  },
};
</script>
