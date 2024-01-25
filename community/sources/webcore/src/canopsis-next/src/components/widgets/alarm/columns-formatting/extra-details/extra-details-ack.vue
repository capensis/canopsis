<template>
  <div>
    <v-tooltip
      class="c-extra-details"
      top
    >
      <template #activator="{ on }">
        <span
          class="c-extra-details__badge purple"
          v-on="on"
        >
          <v-icon
            color="white"
            small
          >
            {{ icon }}
          </v-icon>
        </span>
      </template>
      <div class="text-md-center">
        <strong>{{ $t('alarm.actions.iconsTitles.ack') }}</strong>
        <div>{{ $t('common.by') }} : {{ ack.a }}</div>
        <div>{{ $t('common.date') }} : {{ date }}</div>
        <div v-if="ack.initiator">
          {{ $t('common.initiator') }} : {{ ack.initiator }}
        </div>
        <div
          class="c-extra-details__message"
          v-if="ack.m"
        >
          {{ $tc('common.comment') }} : {{ ack.m }}
        </div>
      </div>
    </v-tooltip>
  </div>
</template>

<script>
import { EVENT_ENTITY_TYPES } from '@/constants';

import { getEntityEventIcon } from '@/helpers/entities/entity/icons';
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';

export default {
  props: {
    ack: {
      type: Object,
      required: true,
    },
  },
  computed: {
    date() {
      return convertDateToStringWithFormatForToday(this.ack.t);
    },

    icon() {
      return getEntityEventIcon(EVENT_ENTITY_TYPES.ack);
    },
  },
};
</script>
