<template lang="pug">
  v-dialog(:value="opened", @input="hideModal", max-width="700", lazy)
    v-form(@submit.prevent="submit")
      v-card
        v-card-title
          span.headline {{ $t('modals.addChangeStateEvent.title') }}
        v-card-text
          v-layout(row)
            v-text-field(
            :label="$t('modals.addChangeStateEvent.output')",
            :error-messages="errors.collect('output')",
            v-model="form.output",
            v-validate="'required'",
            data-vv-name="output"
            )
          v-layout(row)
            date-time-picker(v-model="form.dateTimeStart")
          v-layout(row)
            date-time-picker(v-model="form.dateTimeEnd")
          v-layout(row)
            v-tabs(v-model="activeTab")
              v-tab(href="#simple") Simple
              v-tab(href="#advanced") Advanced
              v-tab(href="#text-input") Text input
              v-tab-item(id="simple")
                div
                  div
                    v-select(
                    :items="items",
                    label="Select"
                    )
                  div
                    date-time-picker(v-model="form.dateTimeEnd")
                  div
                    v-select(
                    label="Select",
                    :items="itemsSecond",
                    multiple,
                    chips
                    )
                  div
                    v-text-field(type="number")
                  div
                    v-text-field(type="number")
              v-tab-item(id="advanced")
                div Advanced
              v-tab-item(id="text-input")
                div Text input
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
      activeTab: 'simple',
      form: {
        name: '',
        dateTimeStart: now,
        dateTimeEnd: now,
      },
      items: [
        'Secondly',
        'Minutely',
        'Hourly',
        'Daily',
        'Weekly',
        'Monthly',
        'Yearly',
      ],
      itemsSecond: [
        'Monday',
        'Tuesday',
        'Wednesday',
        'Thursday',
        'Friday',
        'Saturday',
        'Sunday',
      ],
    };
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
