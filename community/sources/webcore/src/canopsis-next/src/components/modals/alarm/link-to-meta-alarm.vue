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
import { ref, computed } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { mapIds } from '@/helpers/array';
import { isAlarmStateNotOk } from '@/helpers/entities/alarm/form';
import { metaAlarmLinkToForm, formToMetaAlarmLinkRequest } from '@/helpers/entities/meta-alarm/link/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

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
  props: {
    modal: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const form = ref(metaAlarmLinkToForm());

    const { modals, config, close } = useInnerModal(props);

    const alarms = computed(() => config.value.items?.filter(isAlarmStateNotOk) ?? []);

    const submitMethod = async () => {
      const data = {
        ...formToMetaAlarmLinkRequest(form.value),
        alarms: mapIds(alarms.value),
      };

      await config.value?.action?.(data);

      modals.hide();
    };

    const { submit, submitting, isDisabled } = useSubmittableForm({
      form,
      method: submitMethod,
    });

    useFormConfirmableCloseModal({
      form,
      submit: submitMethod,
      close,
    });

    return {
      form,
      alarms,
      submit,
      submitting,
      isDisabled,
    };
  },
};
</script>
