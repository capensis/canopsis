<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title
        span.headline {{ $t(title) }}
      v-card-text
        v-container
          v-layout(row)
            v-flex.text-xs-center
              alarm-general-table
          v-layout(row)
            v-divider.my-3
          v-layout(row)
            v-text-field(
            :label="$t('modals.createCancelEvent.output')",
            :error-messages="errors.collect('output')",
            v-model="form.output",
            v-validate="'required'",
            data-vv-name="output"
            )
      v-card-actions
        v-btn(type="submit", :disabled="errors.any()", color="primary") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import AlarmGeneralTable from '@/components/tables/alarm-general.vue';
import EventEntityMixin from '@/mixins/event-entity';
import ModalMixin from '@/mixins/modal';
import { EVENT_TYPES } from '@/config';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    AlarmGeneralTable,
  },
  mixins: [ModalMixin, EventEntityMixin],
  data() {
    return {
      form: {
        output: '',
      },
    };
  },
  computed: {
    item() {
      return this.getItem(this.config.itemType, this.config.itemId);
    },

    title() {
      return this.config.title || 'modals.createCancelEvent.title';
    },

    eventType() {
      return this.config.eventType || EVENT_TYPES.cancel;
    },
  },

  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = { ...this.form };

        if (this.config.eventType === EVENT_TYPES.cancel) {
          data.cancel = 1;
        }

        await this.createEvent(this.config.eventType, this.item, data);

        this.hideModal();
      }
    },
  },
};
</script>
