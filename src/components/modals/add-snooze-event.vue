<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title
        span.headline {{ $t('modals.addSnoozeEvent.title') }}
      v-card-text
        v-container
          v-layout(row)
            v-flex(xs8)
              v-text-field(
              type="number",
              :label="$t('modals.addSnoozeEvent.duration')",
              :error-messages="errors.collect('duration')",
              v-model="form.duration",
              v-validate="'required|numeric|min_value:1'",
              data-vv-name="duration"
              )
            v-flex(xs4)
              v-select(:items="availableTypes", v-model="form.durationType")
                template(slot="selection" slot-scope="data")
                  div.input-group__selections__comma {{ $tc(data.item.text, 2) }}
                template(slot="item" slot-scope="data")
                  div.list__tile__title {{ $tc(data.item.text, 2) }}
      v-card-actions
        v-btn(type="submit", :disabled="errors.any()", color="primary") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import moment from 'moment';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  data() {
    const availableTypes = [
      { key: 'seconds', text: 'common.times.second' },
      { key: 'minutes', text: 'common.times.minute' },
      { key: 'hours', text: 'common.times.hour' },
      { key: 'days', text: 'common.times.day' },
      { key: 'weeks', text: 'common.times.week' },
      { key: 'months', text: 'common.times.month' },
      { key: 'years', text: 'common.times.year' },
    ];

    return {
      form: {
        duration: 1,
        durationType: availableTypes[0],
      },
      availableTypes,
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const duration = moment.duration(
          parseInt(this.form.duration, 10),
          this.form.durationType.key,
        ).asSeconds();

        console.log(duration);
      }
    },
  },
};
</script>
