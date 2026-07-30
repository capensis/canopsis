<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <kpi-filter-form />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
:disabled="isSubmitting"
:loading="isSubmitting"
color="primary"
type="submit">
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { MODALS, PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';

import { useForm } from '@/hooks/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import KpiFilterForm from '@/components/other/kpi/filters/form/kpi-filter-form.vue';

export default {
  name: MODALS.createKpiFilter,

  components: { KpiFilterForm },

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
    const filter = computed(() => config.value.filter ?? {});
    const isNew = computed(() => !filter.value._id);
    const title = computed(() => config.value.title ?? t('modals.createFilter.create.title'));
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

    const form = useForm({
      initialValues: {
        name: filter.value.name ?? '',
        patterns: filterPatternsToForm(filter.value, [PATTERNS_FIELDS.entity]),
      },
    });
    const { handleSubmit, isSubmitting, values } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.({
        name: values.value.name,
        ...formFilterToPatterns(values.value.patterns, [PATTERNS_FIELDS.entity]),
      });

      await config.value.afterSubmit?.(result);

      closeModal();

      return result;
    });

    useFormConfirmableCloseModal({ form: values, submit, close: closeModal });

    return {
      t,
      title,
      submitLabel,
      isSubmitting,
      closeModal,
      submit,
    };
  },
};
</script>
