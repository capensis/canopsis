<template lang="pug">
  div
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('actions.title') }}
    div.white
      v-layout(row, wrap)
        v-flex(xs4)
          search-field(v-model="searchingText")
        v-flex(v-show="selected.length", xs4)
          v-btn(icon, @click="")
            v-icon(color="error") delete
      v-data-table(
      v-model="selected",
      :items="actions",
      :loading="actionsPending",
      :headers="headers",
      :search="searchingText",
      item-key="_id",
      select-all,
      expand
      )
        template(slot="progress")
          v-fade-transition
            v-progress-linear(height="2", indeterminate, color="primary")
        template(slot="items", slot-scope="props")
          tr(@click="props.expanded = !props.expanded")
            td(@click.stop="")
              v-checkbox(v-model="props.selected", primary, hide-details)
            td {{ props.item._id }}
            td {{ props.item.type }}
            td
              v-layout
                v-flex
                  v-btn(icon, small, @click.stop="")
                    v-icon edit
                  v-btn.error--text(icon, small, @click.stop="showDeleteActionModal(props.item._id)")
                    v-icon(color="error") delete
        template(slot="expand", slot-scope="{ item }")
          actions-list-expand-item(:action="item")
</template>

<script>
import Ellipsis from '@/components/tables/ellipsis.vue';
import SearchField from '@/components/forms/fields/search-field.vue';
import RefreshBtn from '@/components/other/view/refresh-btn.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';

import entitiesActionMixin from '@/mixins/entities/action';

import ActionsListExpandItem from './actions-list-expand-item.vue';

export default {
  components: {
    Ellipsis,
    SearchField,
    RefreshBtn,
    RecordsPerPage,
    ActionsListExpandItem,
  },
  mixins: [entitiesActionMixin],
  props: {
    showDeleteActionModal: {
      type: Function,
      default: () => () => {},
    },
  },
  data() {
    return {
      searchingText: '',
      selected: [],
    };
  },
  computed: {
    headers() {
      return [
        { text: this.$t('actions.table.id'), value: 'id' },
        { text: this.$t('actions.table.type'), value: 'type' },
        { text: this.$t('common.actionsLabel'), value: 'actions' },
      ];
    },
  },
  watch: {
    query() {
      this.selected = [];
    },
  },
};
</script>
