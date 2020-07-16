<template lang="pug">
  modal-wrapper
    template(slot="title")
      span {{ $t('modals.selectExceptionsDatesLists.title') }}
    template(slot="text")
      v-layout(row)
        search-field(@submit="search")
      v-data-table(
        v-model="selected",
        :headers="headers",
        :items="items",
        :loading="pending",
        :total-items="items.length",
        :pagination.sync="pagination",
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
import { createNamespacedHelpers } from 'vuex';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import vuetifyPaginationMixinCreator from '@/mixins/vuetify/pagination-creator';

import SearchField from '@/components/forms/fields/search-field.vue';

import ModalWrapper from '../modal-wrapper.vue';

const { mapActions } = createNamespacedHelpers('pbehaviorException');

export default {
  name: MODALS.selectExceptionsDatesLists,
  components: { ModalWrapper, SearchField },
  mixins: [
    modalInnerMixin,
    vuetifyPaginationMixinCreator({
      mutating: true,
    }),
  ],
  data() {
    return {
      pending: false,
      selected: [],
      query: {
        page: 1,
        limit: PAGINATION_LIMIT,
        search: '',
      },
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
  mounted() {
    // this.fetchList();
  },
  methods: {
    ...mapActions({
      fetchPbehaviorExceptionsListWithoutStore: 'fetchListWithoutStore',
    }),

    fetchList() {
      this.pending = true;

      this.fetchPbehaviorExceptionsListWithoutStore({
        search: this.searchingText,
      });

      this.pending = false;
    },

    search(search) {
      this.query.search = search;
    },
  },
};
</script>
