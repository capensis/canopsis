<template>
  <v-layout align-center>
    <v-flex xs6>
      <c-number-field
        v-field="value.interval"
        :label="$t('common.every')"
        :min="1"
        name="interval"
        required
      />
    </v-flex>
    <v-flex
      class="pl-2 text--secondary"
      xs6
    >
      {{ intervalTimeString }}
    </v-flex>
  </v-layout>
</template>

<script>
import { computed } from 'vue';
import { RRule } from 'rrule';

import { useI18n } from '@/hooks/i18n';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { te, tc } = useI18n();

    const intervalTimeString = computed(() => {
      const timeMessageKey = {
        [RRule.HOURLY]: 'common.times.hour',
        [RRule.DAILY]: 'common.times.day',
        [RRule.WEEKLY]: 'common.times.week',
        [RRule.MONTHLY]: 'common.times.month',
        [RRule.YEARLY]: 'common.times.year',
      }[props.value.freq];

      if (!timeMessageKey || !te(timeMessageKey)) {
        return '';
      }

      return tc(timeMessageKey, props.value.interval || 1);
    });

    return {
      intervalTimeString,
    };
  },
};
</script>
