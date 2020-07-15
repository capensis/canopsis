<template lang="pug">
  modal-wrapper
    template(slot="title")
      span {{ $t('modals.selectExceptionsDatesLists.title') }}
    template(slot="text")
      v-layout(row)
        search-field(v-model="searchingText", @submit="search")
      v-data-table(
        v-model="selected",
        :headers="headers",
        :items="items",
        :loading="pending",
        :total-items="items.length",
        item-key="id",
        select-all
      )
        template(slot="items", slot-scope="props")
          tr
            td
              v-checkbox-functional(v-model="props.selected", primary, hide-details)
            td {{ props.item.name }}
    template(slot="actions")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn.primary(type="submit") {{ $t('common.submit') }}
    div
</template>

<script>
import { MODALS } from '@/constants';

import SearchField from '@/components/forms/fields/search-field.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.selectExceptionsDatesLists,
  components: { ModalWrapper, SearchField },
  data() {
    return {
      pending: false,
      selected: [],
      searchingText: '',
    };
  },
  computed: {
    headers() {
      return [{ text: 'name', value: 'name' }];
    },
    items() {
      return [
        { id: 1, name: 'root' },
        { id: 2, name: 'canopsis' },
      ];
    },
  },
  methods: {
    search() {
    },
  },
};
</script>
