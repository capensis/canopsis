<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <filter-form v-bind="patternsProps" />
      </template>

      <template #actions>
        <v-btn :disabled="isSubmitting" variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="isSubmitting"
          :loading="isSubmitting"
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
import { omit } from 'lodash';
import { computed } from 'vue';

import { MODALS, PATTERNS_FIELDS } from '@/constants';

import { filterToForm, formToFilter } from '@/helpers/entities/filter/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import FilterForm from '@/components/other/filter/form/filter-form.vue';

export default {
  name: MODALS.createFilter,

  components: { FilterForm },

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
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.save'));

    const patternsFields = computed(() => {
      const { withAlarm, withEntity, withPbehavior, withEvent, withServiceWeather } = config.value;

      return [
        withAlarm && PATTERNS_FIELDS.alarm,
        withEntity && PATTERNS_FIELDS.entity,
        withPbehavior && PATTERNS_FIELDS.pbehavior,
        withEvent && PATTERNS_FIELDS.event,
        withServiceWeather && PATTERNS_FIELDS.serviceWeather,
      ].filter(Boolean);
    });

    const form = useForm({
      initialValues: filterToForm(filter.value, patternsFields.value),
    });
    const { handleSubmit, isSubmitting, values } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });
    const patternsProps = computed(() =>
      omit(config.value, [
        'title',
        'action',
        'afterSubmit',
        'filter',
        'corporate',
        'widgetType',
        'withoutLink',
        'withTitle',
      ]),
    );

    const submit = handleSubmit(async () => {
      const newFilter = formToFilter(values.value, patternsFields.value, config.value.corporate);
      const result = await config.value.action?.(newFilter);

      await config.value.afterSubmit?.(result);

      closeModal();

      return result;
    });

    useFormConfirmableCloseModal({ form: values, submit, close: closeModal });

    return {
      t,
      title,
      patternsProps,
      isSubmitting,
      submitLabel,
      submit,
      closeModal,
    };
  },
};
</script>
