<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ $t('common.filters') }}
    template(#text="")
      filters-list-component(
        :widget-id="config.widgetId",
        :with-pbehavior="config.withPbehavior",
        :with-alarm="config.withAlarm",
        :with-entity="config.withEntity",
        :private="config.private",
        :addable="config.hasAccessToAddFilter",
        :editable="config.hasAccessToEditFilter"
      )
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import FiltersListComponent from '@/components/other/filter/filters-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Confirmation modal
 */
export default {
  name: MODALS.filtersList,
  components: { FiltersListComponent, ModalWrapper },
  mixins: [modalInnerMixin],
};
</script>
