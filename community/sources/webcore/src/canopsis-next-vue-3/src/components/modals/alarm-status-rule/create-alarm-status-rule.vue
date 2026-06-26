<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper
      text-class="position-relative"
      close
    >
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <alarm-status-rule-form
          :disablable="config.disablable"
          :flapping="config.flapping"
          :default-rule="isDefaultRule"
        />
      </template>

      <template #actions>
        <v-btn
          variant="text"
          @click="closeModal"
        >
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="submitting"
          :loading="submitting"
          color="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { DEFAULT_RESOLVE_RULE_ID } from '@/constants';

import { alarmStatusRuleToForm, formToAlarmStatusRule } from '@/helpers/entities/alarm-status-rule/form';

import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useModals } from '@/hooks/modals';

import AlarmStatusRuleForm from '@/components/other/alarm-status-rule/form/alarm-status-rule-form.vue';

export default {
  components: { AlarmStatusRuleForm },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});

    const isNew = computed(() => !config.value.rule?._id || config.value.duplicating);
    const title = computed(() => config.value.title);
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));
    const isDefaultRule = computed(() => config.value.rule?._id === DEFAULT_RESOLVE_RULE_ID);
    const form = useForm({
      initialValues: alarmStatusRuleToForm(config.value.rule, config.value.flapping),
    });
    const {
      values,
      handleSubmit,
      submitting,
    } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToAlarmStatusRule(values.value));

      await config.value.afterSubmit?.(result);

      closeModal();

      return result;
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    return {
      t,
      config,
      title,
      submitLabel,
      submitting,
      isDefaultRule,
      closeModal,
      submit,
    };
  },
};
</script>
