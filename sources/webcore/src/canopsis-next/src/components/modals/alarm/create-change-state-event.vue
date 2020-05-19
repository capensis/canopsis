<template lang="pug">
  v-form(data-test="createChangeStateEventModal", @submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createChangeStateEvent.title') }}
      template(slot="text")
        v-container
          change-state-field(
            v-model="form",
            :label="$t('modals.createChangeStateEvent.fields.output')"
          )
      template(slot="actions")
        v-btn(
          data-test="createChangeStateEventCancelButton",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled",
          data-test="createChangeStateEventSubmitButton",
          type="submit"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS, ENTITIES_STATES, EVENT_ENTITY_TYPES } from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';
import submittableMixin from '@/mixins/submittable';
import entitiesInfoMixin from '@/mixins/entities/info';

import ChangeStateField from '@/components/other/action/form/fields/change-state.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create a 'change-state' event
 */
export default {
  name: MODALS.createChangeStateEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ChangeStateField, ModalWrapper },
  mixins: [entitiesInfoMixin, modalInnerItemsMixin, eventActionsAlarmMixin, submittableMixin()],
  data() {
    return {
      form: {
        output: '',
        state: ENTITIES_STATES.major,
      },
    };
  },
  mounted() {
    this.form.state = this.firstItem.v.state.val;
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.createEvent(EVENT_ENTITY_TYPES.changeState, this.items, this.form);

        this.$modals.hide();
      }
    },
  },
};
</script>
