<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title
        span.headline {{ $t(title) }}
      v-card-text
        v-container
          v-layout(row)
            v-flex.text-xs-center
              alarm-general-table(:items="items")
          v-layout(row)
            v-divider.my-3
          v-layout(row)
            v-text-field(
            :label="$t('modals.createCancelEvent.fields.output')",
            :error-messages="errors.collect('output')",
            v-model="form.output",
            v-validate="'required'",
            data-vv-name="output"
            )
      v-card-actions
        v-btn(type="submit", :disabled="errors.any()", color="primary") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import AlarmGeneralTable from '@/components/tables/alarm/general.vue';
import ModalInnerItemsMixin from '@/mixins/modal/modal-inner-items';
import EventActionsMixin from '@/mixins/event-actions';
import { EVENT_ENTITY_TYPES, MODALS } from '@/constants';

export default {
  name: MODALS.createCancelEvent,

  $_veeValidate: {
    validator: 'new',
  },
  components: {
    AlarmGeneralTable,
  },
  mixins: [ModalInnerItemsMixin, EventActionsMixin],
  data() {
    return {
      form: {
        output: '',
      },
    };
  },
  computed: {
    title() {
      return this.config.title || 'modals.createCancelEvent.title';
    },

    eventType() {
      return this.config.eventType || EVENT_ENTITY_TYPES.cancel;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = { ...this.form };

        if (this.eventType === EVENT_ENTITY_TYPES.cancel) {
          data.cancel = 1;
        }

        await this.createEvent(this.eventType, this.items, data);

        this.hideModal();
      }
    },
  },
};
</script>
