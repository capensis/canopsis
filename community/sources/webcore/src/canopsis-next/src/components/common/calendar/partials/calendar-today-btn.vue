<template>
  <v-tooltip bottom>
    <template #activator="{ on }">
      <v-btn
        class="mx-2 my-1"
        v-on="on"
        :icon="$vuetify.breakpoint.smAndDown"
        depressed
        @click="$emit('click')"
      >
        <span v-if="$vuetify.breakpoint.mdAndUp">{{ $t('calendar.today') }}</span>
        <v-icon v-else>
          today
        </v-icon>
      </v-btn>
    </template>
    <span>{{ todayDate }}</span>
  </v-tooltip>
</template>

<script>
import { convertDateToString, getNowTimestamp } from '@/helpers/date/date';

export default {
  computed: {
    todayDate() {
      return convertDateToString(getNowTimestamp(), 'dddd, MMMM D');
    },
  },
};
</script>
