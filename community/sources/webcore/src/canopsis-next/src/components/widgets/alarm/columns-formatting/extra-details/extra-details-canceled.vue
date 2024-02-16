<template>
  <div>
    <v-tooltip
      class="c-extra-details"
      top
    >
      <template #activator="{ on }">
        <span
          class="c-extra-details__badge blue-grey"
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
        <strong>{{ $t('alarm.actions.iconsTitles.canceled') }}</strong>
        <div>{{ $t('common.by') }} : {{ canceled.a }}</div>
        <div>{{ $t('common.date') }} : {{ date }}</div>
        <div
          class="c-extra-details__message"
          v-if="canceled.m"
        >
          {{ $tc('common.comment') }} : {{ canceled.m }}
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
    canceled: {
      type: Object,
      required: true,
    },
  },
  computed: {
    date() {
      return convertDateToStringWithFormatForToday(this.canceled.t);
    },

    icon() {
      return getEntityEventIcon(EVENT_ENTITY_TYPES.delete);
    },
  },
};
</script>
