<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.createChangeStateEvent.title') }}</span>
      </template>

      <template #text>
        <alarm-general-table
          v-if="config.items"
          :items="config.items"
          class="mb-4"
        />
        <c-change-state-field
          v-model="values"
          :label="t('modals.createChangeStateEvent.fields.output')"
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

import { ALARM_STATES, MODALS } from '@/constants';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';

/**
 * Modal to create a 'change-state' event
 */
export default {
  name: MODALS.createChangeStateEvent,

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
    const items = computed(() => config.value.items ?? []);
    const [firstItem] = items.value;

    const form = useForm({
      initialValues: {
        comment: '',
        state: firstItem && items.value.length === 1 ? firstItem.v.state.val : ALARM_STATES.major,
      },
    });
    const { values } = form;

    const { submitting, submit } = useSubmittableForm({
      method: form.handleSubmit(async () => {
        await config.value.action?.(values.value);
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
      config,
      values,
      submitting,
      submit,
      close,
    };
  },
};
</script>
