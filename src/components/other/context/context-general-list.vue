<template lang="pug">
div
    v-subheader {{ $t('tables.contextList.title') }}
    v-toolbar.toolbar(dense, flat)
      v-text-field(
      label="Search",
      v-model="searchingText",
      hide-details,
      single-line,
      @keyup.enter="submit",
      )
      v-btn(icon, @click="submit")
        v-icon search
      v-btn(icon, @click="$emit('update:selectedIds',selectedEntities)")
        v-icon done
    v-data-table(
      :no-data-text="this.$t('tables.contextList.noDataText')",
      :headers="headers",
      :items="contextEntities",
      :loading="pending",
      v-model="selectedEntities",
      select-all,
      item-key="_id",
      )
      template(slot="items", slot-scope="props")
        td
          v-checkbox(
            v-model="props.selected",
            :value="props._id",
            primary,
            hide-details,
          )
        td.text-xs-left {{ props.item.name }}
        td.text-xs-left {{ props.item._id}}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('context');
export default {
  data() {
    return {
      searchingText: '',
      selectedEntities: [],
      headers: [
        {
          text: this.$t('tables.contextList.name'),
          sortable: false,
        },
        {
          text: this.$t('tables.contextList.id'),
          sortable: false,
        },
      ],
    };
  },
  computed: {
    ...mapGetters({
      contextEntities: 'itemsGeneralList',
      pending: 'pendingGeneralList',
    }),
  },
  methods: {
    ...mapActions({
      fetchContextEntities: 'fetchGeneralList',
    }),
    submit() {
      this.fetchContextEntities({
        params: {
          _filter: this.$options.filters.formatContextSearch(this.searchingText),
        },
      });
    },
  },
};
</script>

