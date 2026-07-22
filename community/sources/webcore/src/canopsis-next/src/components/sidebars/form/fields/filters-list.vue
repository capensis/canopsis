<template>
  <filters-list
    v-field="filters"
    :addable="addable"
    :editable="editable"
    :name="name"
    :required="required"
    @add="showCreateFilterModal"
    @edit="showEditFilterModal"
    @delete="showDeleteFilterModal"
  />
</template>

<script>
import { pick } from 'lodash';
import { computed, inject } from 'vue';

import { MODALS } from '@/constants';

import { uuid } from '@/helpers/uuid';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useArrayModelField } from '@/hooks/form/array-model-field';

import FiltersList from '@/components/other/filter/filters-list.vue';

export default {
  inject: ['$sidebar'],
  components: { FiltersList },
  model: {
    prop: 'filters',
    event: 'input',
  },
  props: {
    widgetId: {
      type: String,
      required: false,
    },
    filters: {
      type: Array,
      default: () => [],
    },
    addable: {
      type: Boolean,
      default: false,
    },
    editable: {
      type: Boolean,
      default: false,
    },
    withAlarm: {
      type: Boolean,
      default: false,
    },
    withEntity: {
      type: Boolean,
      default: false,
    },
    withPbehavior: {
      type: Boolean,
      default: false,
    },
    withServiceWeather: {
      type: Boolean,
      default: false,
    },
    entityTypes: {
      type: Array,
      required: false,
    },
    entityCountersType: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'filters',
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const sidebar = inject('$sidebar');

    const { t } = useI18n();
    const modals = useModals();

    const {
      addItemIntoArray,
      updateItemInArray,
      removeItemFromArray,
    } = useArrayModelField(props, emit);

    const modalConfig = computed(() => ({
      ...pick(props, [
        'withAlarm',
        'withEntity',
        'withPbehavior',
        'withServiceWeather',
        'entityTypes',
        'entityCountersType',
      ]),

      widgetType: sidebar?.config?.widget?.type,
      withTitle: true,
      withoutLink: true,
    }));

    const showCreateFilterModal = () => {
      modals.show({
        name: MODALS.createFilter,
        config: {
          ...modalConfig.value,

          title: t('modals.createFilter.create.title'),
          action: newFilter => addItemIntoArray({
            ...newFilter,

            _id: uuid('filter'),
            widget: props.widgetId,
            is_user_preference: false,
          }),
        },
      });
    };

    const showEditFilterModal = (filter, index) => {
      modals.show({
        name: MODALS.createFilter,
        config: {
          ...modalConfig.value,

          filter,
          title: t('modals.createFilter.edit.title'),
          action: newFilter => updateItemInArray(index, {
            ...newFilter,

            widget: props.widgetId,
            _id: filter._id,
          }),
        },
      });
    };

    const showDeleteFilterModal = (filter, index) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => removeItemFromArray(index),
        },
      });
    };

    return {
      showCreateFilterModal,
      showEditFilterModal,
      showDeleteFilterModal,
    };
  },
};
</script>
