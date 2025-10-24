<template>
  <v-layout class="gap-2" column>
    <span class="grey--text">{{ $t('storageSetting.entityUnlinked.archiveUnlinkedAfter') }}</span>
    <c-duration-field
      v-field="duration"
      :units-label="$t('common.unit')"
      :units="timeUnits"
      required
    />
  </v-layout>
</template>

<script>
import { toRef } from 'vue';

import { TIME_UNITS } from '@/constants';

import { useEnabledDurationField } from '@/components/forms/fields/duration/hooks/enabled-duration-field';

export default {
  model: {
    prop: 'duration',
    event: 'input',
  },
  props: {
    duration: {
      type: Object,
      default: () => ({
        value: 60,
        unit: TIME_UNITS.day,
      }),
    },
    name: {
      type: String,
      default: 'duration',
    },
  },
  setup(props) {
    const {
      timeUnits,
    } = useEnabledDurationField({
      duration: toRef(props, 'duration'),
      name: toRef(props, 'name'),
    });

    return {
      timeUnits,
    };
  },
};
</script>
