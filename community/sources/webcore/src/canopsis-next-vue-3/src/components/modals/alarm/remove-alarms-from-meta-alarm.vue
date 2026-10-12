<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ config.title }}</span>
      </template>

      <template #text>
        <alarm-general-table :items="config.items" class="mb-4" />
        <c-name-field
          v-model="comment"
          :label="t('common.comment')"
          :required="isCommentRequired"
          name="comment"
          autofocus
        />
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

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';

/**
 * Modal to remove alarms from meta alarm
 */
export default {
  name: MODALS.removeAlarmsFromMetaAlarm,

  components: { AlarmGeneralTable },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = useForm({
      initialValues: {
        comment: '',
      },
    });
    const comment = form.model('comment');
    const { handleSubmit, isSubmitting: submitting } = form;

    const isCommentRequired = computed(() => config.value.isCommentRequired ?? true);

    const submit = handleSubmit(async () => {
      await config.value?.action?.({
        ...form.values.value,
        alarms: mapIds(config.value.items),
      });

      close();
    });

    useFormConfirmableCloseModal({
      form: form.values,
      submit,
      close,
    });

    return {
      t,
      config,
      comment,
      isCommentRequired,
      submitting,
      submit,
      close,
    };
  },
};
</script>
