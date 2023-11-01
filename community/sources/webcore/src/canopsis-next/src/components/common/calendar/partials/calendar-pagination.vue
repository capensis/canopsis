<template>
  <div>
    <v-tooltip bottom>
      <template #activator="{ on }">
        <v-btn
          class="mx-2"
          v-on="on"
          icon
          depressed
          @click="$emit('prev')"
        >
          <v-icon>keyboard_arrow_left</v-icon>
        </v-btn>
      </template>
      <span>{{ prevLabel }}</span>
    </v-tooltip>
    <calendar-period-picker
      :type="type"
      :focus="focus"
      @input="selectPeriod"
    />
    <v-tooltip bottom>
      <template #activator="{ on }">
        <v-btn
          class="mx-2"
          v-on="on"
          icon
          depressed
          @click="$emit('next')"
        >
          <v-icon>keyboard_arrow_right</v-icon>
        </v-btn>
      </template>
      <span>{{ nextLabel }}</span>
    </v-tooltip>
  </div>
</template>

<script>
import { lowerCase } from 'lodash';

import CalendarPeriodPicker from './calendar-period-picker.vue';

export default {
  components: { CalendarPeriodPicker },
  props: {
    focus: {
      type: Date,
      required: false,
    },
    type: {
      type: String,
      required: false,
    },
  },
  computed: {
    typeString() {
      return lowerCase(this.$t(`calendar.${this.type}`));
    },

    prevLabel() {
      return [
        this.$t('common.previous'),
        this.typeString,
      ].join(' ');
    },

    nextLabel() {
      return [
        this.$t('common.next'),
        this.typeString,
      ].join(' ');
    },
  },
  methods: {
    selectPeriod(value) {
      this.$emit('update:focus', value);
    },
  },
};
</script>
