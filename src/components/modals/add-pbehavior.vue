<template lang="pug">
  v-dialog(:value="opened", @input="hideModal", max-width="700")
    v-form(@submit.prevent="submit")
      v-card
        v-card-title
          span.headline {{ $t('modals.addChangeStateEvent.title') }}
        v-card-text
          v-container
            v-layout(row)
              v-text-field(
              :label="$t('modals.addChangeStateEvent.output')",
              :error-messages="errors.collect('output')",
              v-model="form.output",
              v-validate="'required'",
              data-vv-name="output"
              )
            v-layout(row)
              date-time-picker(v-model="form.dateTime")

        v-card-actions
          v-btn(type="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import moment from 'moment';

import DateTimePicker from '@/components/forms/date-time-picker.vue';

import ModalMixin from './modal-mixin';

export default {
  name: 'add-pbehavior',
  mixins: [ModalMixin],
  components: { DateTimePicker },
  data() {
    const now = moment();

    return {
      menu: false,
      form: {
        name: '',
        dateTime: now,
      },
      value: 'date',
    };
  },
  computed: {
    fullDate() {
      return `${this.form.date} ${this.form.time}`;
    },
  },
  methods: {
    ok() {
      this.$refs.menu.save();

      setTimeout(() => {
        this.value = 'date';
      }, 300);
    },
  },
};
</script>

<style>
  .tabs__container--centered .tabs__div,
  .tabs__container--fixed-tabs .tabs__div,
  .tabs__container--icons-and-text .tabs__div {
    min-width: 140px;
  }

  .menu__content {
    max-width: 100%;
  }

  .dropdown-footer, .menu__content {
    background-color: #fff;
  }

  .date-picker-table {
    height: 246px;
  }
</style>
