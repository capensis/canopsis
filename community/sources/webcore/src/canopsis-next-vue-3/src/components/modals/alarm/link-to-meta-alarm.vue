<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        {{ t('modals.linkToMetaAlarm.title') }}
      </template>

      <template #text>
        <c-layout gap="2" column>
          <alarm-general-table :items="alarms" class="mb-4" />
          <link-meta-alarm-form :form="values" @update:form="updateForm" />
        </c-layout>
      </template>

      <template #actions>
        <v-btn :disabled="submitting" variant="text" @click="close">
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn
          :loading="submitting"
          :disabled="submitting"
          color="primary"
          type="submit"
        >
          {{ t('common.saveChanges') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { mapIds } from '@/helpers/array';
import { isAlarmStateNotOk } from '@/helpers/entities/alarm/form';
import { formToMetaAlarmLinkRequest, metaAlarmLinkToForm } from '@/helpers/entities/meta-alarm/link/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';
import LinkMetaAlarmForm from '@/components/widgets/alarm/forms/link-meta-alarm-form.vue';

/**
 * Modal to manage alarms in manual meta alarm
 */
export default {
  name: MODALS.linkToMetaAlarm,

  components: {
    AlarmGeneralTable,
    LinkMetaAlarmForm,
  },

  props: {
    modal: {
      type: Object,
      default: () => ({}),
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = useForm({
      initialValues: metaAlarmLinkToForm(),
    });
    const { values } = form;

    const alarms = computed(() => config.value.items?.filter(isAlarmStateNotOk) ?? []);

    const updateForm = (nextForm) => {
      values.value = nextForm;
    };

    const { submit, submitting } = useSubmittableForm({
      method: form.handleSubmit(async () => {
        const data = {
          ...formToMetaAlarmLinkRequest(values.value),
          alarms: mapIds(alarms.value),
        };

        await config.value?.action?.(data);

        close();
      }),
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close,
    });

    return {
      t,
      form,
      values,
      alarms,
      submit,
      submitting,
      close,
      updateForm,
    };
  },
};
</script>
