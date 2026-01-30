<template>
  <v-chip
    :close="closable"
    :color="color"
    class="c-advanced-search__chip c-advanced-search__array-chip"
    @click.prevent=""
    @click:close="close"
  >
    <date-time-picker-menu
      v-for="chip in chips"
      v-field="value[chip]"
      :key="chip"
    >
      <template #activator="{ on, value: chipValue }">
        <v-chip v-on="on">
          <span class="font-italic grey--text">
            {{ $t(`common.${chip}`) }}:
          </span> {{ chipValue | date('long', '--.--.---- --:--') }}
        </v-chip>
      </template>
    </date-time-picker-menu>
  </v-chip>
</template>

<script>
import DateTimePickerMenu from '@/components/forms/fields/date-time-picker/date-time-picker-menu.vue';

export default {
  components: {
    DateTimePickerMenu,
  },
  props: {
    value: {
      type: Object,
      default: () => ({ from: 0, to: 0 }),
    },
    closable: {
      type: Boolean,
      default: false,
    },
    color: {
      type: String,
      required: false,
    },
  },
  setup(props, { emit }) {
    const chips = ['from', 'to'];

    const close = () => emit('close');

    return {
      chips,
      close,
    };
  },
};
</script>
