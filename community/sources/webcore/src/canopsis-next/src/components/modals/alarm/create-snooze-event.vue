<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createSnoozeEvent.title') }}</span>
      </template>
      <template #text="">
        <snooze-event-form
          v-model="form"
          :is-note-required="config.isNoteRequired"
        />
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :loading="submitting"
          :disabled="isDisabled"
          class="primary"
          type="submit"
        >
          {{ $t('common.saveChanges') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS } from '@/constants';

import { snoozeToForm } from '@/helpers/entities/alarm/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import SnoozeEventForm from '@/components/widgets/alarm/forms/snooze-event-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to put a snooze on an alarm
 */
export default {
  name: MODALS.createSnoozeEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    ModalWrapper,
    SnoozeEventForm,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: snoozeToForm(),
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config?.action?.(this.form);

        this.$modals.hide();
      }
    },
  },
};
</script>
