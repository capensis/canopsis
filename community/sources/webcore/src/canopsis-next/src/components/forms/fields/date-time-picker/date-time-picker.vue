<template lang="pug">
  div
    div.v-picker__title.primary.text-xs-center(v-if="label")
      span.headline {{ label }}
    div.date-time-picker__body
      v-layout.py-2(row, align-center, justify-center)
        v-flex.v-date-time-picker__subtitle-wrapper
          span.v-date-time-picker__subtitle(
            :class="{ 'grey--text darken-1': !localValue }"
          ) {{ localValue | date(dateFormat, '−−/−−/−−−−') }}
        v-flex.v-date-time-picker__subtitle-wrapper
          time-picker-field.v-date-time-picker__subtitle(
            :value="localValue | date('timePicker', null)",
            :round-hours="roundHours",
            @input="updateTime"
          )
      div
        v-date-picker(
          :locale="$i18n.locale",
          :value="localValue | date('datePicker', null)",
          color="primary",
          no-title,
          @input="updateDate"
        )
    slot(name="footer", @submit="submit")
      v-divider
      v-layout.mt-1(justify-space-around)
        v-btn(depressed, flat, @click="$listeners.close") {{ $t('common.cancel') }}
        v-btn.primary(@click="submit") {{ $t('common.apply') }}
</template>

<script>
import { isDate } from 'lodash';

import { updateTime, updateDate } from '@/helpers/date/date-time-picker';

import { formBaseMixin } from '@/mixins/form';

import TimePickerField from '../time-picker/time-picker-field.vue';

export default {
  components: { TimePickerField },
  mixins: [formBaseMixin],
  props: {
    value: {
      type: [Date, Number],
      default: () => new Date(),
    },
    label: {
      type: String,
      default: '',
    },
    dateFormat: {
      type: String,
      default: 'short',
    },
    roundHours: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    const milliseconds = isDate(this.value) ? this.value.getTime() : this.value;

    return {
      localValue: milliseconds ? new Date(milliseconds) : null,
    };
  },
  methods: {
    updateTime(time) {
      this.localValue = updateTime(this.localValue, time);
    },

    updateDate(date) {
      this.localValue = updateDate(this.localValue, date);
    },

    submit() {
      this.updateModel(this.localValue);

      this.$emit('close');
    },
  },
};
</script>

<style lang="scss">
  .date-time-picker {
    .date-time-picker__body {
      position: relative;
      width: 290px;
      height: 352px;
      z-index: inherit;
    }

    .v-date-time-picker__subtitle {
      margin-top: -12px;
      line-height: 30px;
      font-size: 18px;
      font-weight: 400;

      &-wrapper {
        text-align: center;
      }
    }

    .v-menu__content {
      max-width: 100%;
    }

    .v-dropdown-footer, &.v-menu__content, .v-tabs__items {
      background-color: #fff;
    }

    .v-date-picker-table {
      height: 260px;
    }

    .v-card {
      box-shadow: none;
    }

    .v-date-picker-table--date .v-btn {
      height: 35px;
      width: 35px;
    }
  }
</style>
