<template>
  <c-no-events-icon
    v-if="isNoEventsStatus"
    :value="idleSince"
    :size="iconSize"
    color="error"
    top
  />
  <c-simple-tooltip
    v-else
    :content="tooltipContent"
    top
    v-on="$listeners"
  >
    <template #activator="{ on }">
      <v-chip
        v-if="hasUpstream"
        :color="iconStyle.color"
        :small="small"
        class="px-2"
        v-on="on"
      >
        <v-icon
          size="16"
          color="white"
          class="mr-2"
        >
          {{ status.icon }}
        </v-icon>
        <v-icon
          size="16"
          color="white"
        >
          $vuetify.icons.flow
        </v-icon>
      </v-chip>
      <v-icon
        v-else
        :size="iconSize"
        :style="iconStyle"
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

export default {
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

    const statusValue = computed(() => 6); // TODO: revert change
    const isNoEventsStatus = computed(() => true);
    const isOngoingStatus = computed(() => statusValue.value === ALARM_STATUSES.ongoing);
    const hasUpstream = computed(() => true); // TODO: change to props.alarm.entity.upstream
    const idleSince = computed(() => 23);
    const resolved = computed(() => !!props.alarm.v.resolved);
    const status = computed(() => formatAlarmStatus(statusValue.value, resolved.value));
    const state = computed(() => formatAlarmState(props.alarm.v.state.val));
    const statusColor = computed(() => (isOngoingStatus.value ? state.value.color : status.value.color));
    const iconSize = computed(() => (props.small ? 24 : undefined));
    const iconStyle = computed(() => ({ color: statusColor.value, caretColor: statusColor.value }));
    const tooltipContent = computed(() => (resolved.value && te(`common.statusResolvedTypes.${statusValue.value}`)
      ? t(`common.statusResolvedTypes.${statusValue.value}`)
      : t(`common.statusTypes.${statusValue.value}`)));

    return {
      statusValue,
      isNoEventsStatus,
      isOngoingStatus,
      hasUpstream,
      idleSince,
      status,
      state,
      iconSize,
      iconStyle,
      tooltipContent,
    };
  },
};
</script>
