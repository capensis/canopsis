<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createCommentEvent.title') }}</span>
      </template>
      <template #text="">
        <v-layout column>
          <template v-if="items.length">
            <alarm-general-table :items="items" />
            <v-divider class="my-3" />
          </template>
          <c-name-field
            v-model="form.comment"
            :label="$tc('common.comment')"
            name="comment"
            required
          />
        </v-layout>
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
          class="primary"
          :loading="submitting"
          :disabled="isDisabled"
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

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create comment an alarm
 */
export default {
  name: MODALS.createCommentEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AlarmGeneralTable, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: {
        comment: '',
      },
    };
  },
  computed: {
    items() {
      return this.config.items ?? [];
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config.action(this.form);

        this.$modals.hide();
      }
    },
  },
};
</script>
