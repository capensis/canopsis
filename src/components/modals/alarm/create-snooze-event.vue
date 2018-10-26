<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title
        span.headline {{ $t('modals.createSnoozeEvent.title') }}
      v-card-text
        v-container
          v-layout(row)
            v-flex(xs8)
              v-text-field(
              type="number",
              :label="$t('modals.createSnoozeEvent.fields.duration')",
              :error-messages="errors.collect('duration')",
              v-model="form.duration",
              v-validate="'required|numeric|min_value:1'",
              data-vv-name="duration"
              )
            v-flex(xs4)
              v-select(:items="availableTypes", v-model="form.durationType", item-value="key")
                template(slot="selection" slot-scope="data")
                  div.input-group__selections__comma {{ $tc(data.item.text, 2) }}
                template(slot="item" slot-scope="data")
                  div.list__tile__title {{ $tc(data.item.text, 2) }}
      v-card-actions
        v-btn(type="submit", :disabled="errors.any()", color="primary") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import moment from 'moment';

import modalInnerItemsMixin from '@/mixins/modal/modal-inner-items';
import eventActionsMixin from '@/mixins/event-actions';
import { MODALS } from '@/constants';

/**
 * Modal to put a snooze on an alarm
 */
export default {
  name: MODALS.createSnoozeEvent,

  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerItemsMixin, eventActionsMixin],
  data() {
    const availableTypes = [
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
        durationType: availableTypes[0].key,
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
          this.form.durationType,
        ).asSeconds();

        await this.createEvent(this.$constants.EVENT_ENTITY_TYPES.snooze, this.items, { duration });

        this.hideModal();
      }
    },
  },
};
</script>
