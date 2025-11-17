<template>
  <alarm-status-chip-with-relations
    v-if="isNoEventsStatus && hasRelations"
    :small="small"
    :alarm="alarm"
    color="error"
    icon-color="error"
    outlined
  >
    <c-no-events-icon
      :value="idleSince"
      size="16"
      color="error"
      class="alarm-column-value-status__icon mr-2"
      top
    />
  </alarm-status-chip-with-relations>
  <c-no-events-icon
    v-else-if="isNoEventsStatus"
    :value="idleSince"
    :size="iconSize"
    color="error"
    class="alarm-column-value-status__icon"
    top
  />
  <alarm-status-chip-with-relations
    v-else-if="hasRelations"
    :small="small"
    :alarm="alarm"
    :style="chipStyle"
  >
    <c-simple-tooltip
      :content="tooltipContent"
      top
      v-on="$listeners"
    >
      <template #activator="{ on }">
        <v-icon
          size="16"
          color="white"
          class="alarm-column-value-status__icon mr-2"
          v-on="on"
        >
          {{ status.icon }}
        </v-icon>
      </template>
    </c-simple-tooltip>
  </alarm-status-chip-with-relations>
  <c-simple-tooltip
    v-else
    :content="tooltipContent"
    top
    v-on="$listeners"
  >
    <template #activator="{ on }">
      <v-icon
        :style="iconStyle"
        :size="iconSize"
        class="alarm-column-value-status__icon"
        v-on="on"
      >
        {{ status.icon }}
      </v-icon>
    </template>
  </c-simple-tooltip>
</template>

<script>
import { computed } from 'vue';

import { ALARM_STATUSES } from '@/constants';

import { formatAlarmState, formatAlarmStatus } from '@/helpers/entities/alarm/formatting';

import { useI18n } from '@/hooks/i18n';

import AlarmStatusChipWithRelations from '../partials/alarm-status-chip-with-relations.vue';

export default {
  components: {
    AlarmStatusChipWithRelations,
  },
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    small: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t, te } = useI18n();

    const statusValue = computed(() => props.alarm.v.status?.val);
    const isNoEventsStatus = computed(() => statusValue.value === ALARM_STATUSES.noEvents);
    const isOngoingStatus = computed(() => statusValue.value === ALARM_STATUSES.ongoing);
    const idleSince = computed(() => props.alarm.entity?.idle_since);
    const resolved = computed(() => !!props.alarm.v.resolved);
    const status = computed(() => formatAlarmStatus(statusValue.value, resolved.value));
    const state = computed(() => formatAlarmState(props.alarm.v.state.val));
    const statusColor = computed(() => (isOngoingStatus.value ? state.value.color : status.value.color));
    const iconSize = computed(() => (props.small ? 24 : undefined));
    const iconStyle = computed(() => ({ color: statusColor.value, caretColor: statusColor.value }));
    const chipStyle = computed(() => ({ backgroundColor: iconStyle.value.color }));
    const tooltipContent = computed(() => (resolved.value && te(`common.statusResolvedTypes.${statusValue.value}`)
      ? t(`common.statusResolvedTypes.${statusValue.value}`)
      : t(`common.statusTypes.${statusValue.value}`)));

    const hasRelations = computed(() => (
      !!props.alarm.entity?.upstream || props.alarm.entity?.downstream_count > 0
    ));

    return {
      statusValue,
      isNoEventsStatus,
      isOngoingStatus,
      hasRelations,
      idleSince,
      status,
      state,
      iconSize,
      iconStyle,
      chipStyle,
      tooltipContent,
    };
  },
};
</script>
