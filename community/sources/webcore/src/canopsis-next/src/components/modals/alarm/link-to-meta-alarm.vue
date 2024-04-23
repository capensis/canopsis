<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ $t('modals.linkToMetaAlarm.title') }}
      </template>
      <template #text="">
        <v-layout column>
          <alarm-general-table :items="alarms" class="mb-4" />
          <link-meta-alarm-form v-model="form" />
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
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { mapIds } from '@/helpers/array';
import { isAlarmStateNotOk } from '@/helpers/entities/alarm/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';
import LinkMetaAlarmForm from '@/components/widgets/alarm/forms/link-meta-alarm-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to manage alarms in manual meta alarm
 */
export default {
  name: MODALS.linkToMetaAlarm,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    AlarmGeneralTable,
    LinkMetaAlarmForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: {
        metaAlarm: null,
        comment: '',
        auto_resolve: false,
      },
    };
  },
  computed: {
    alarms() {
      return this.config.items?.filter(isAlarmStateNotOk) ?? [];
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = {
          comment: this.form.comment,
          auto_resolve: this.form.auto_resolve,
          alarms: mapIds(this.alarms),
        };

        if (this.form.metaAlarm?._id) {
          data.id = this.form.metaAlarm?._id;
        } else {
          data.name = this.form.metaAlarm;
        }

        await this.config?.action?.(data);

        this.$modals.hide();
      }
    },
  },
};
</script>
