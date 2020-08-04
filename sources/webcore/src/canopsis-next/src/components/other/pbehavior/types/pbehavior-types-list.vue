<template lang="pug">
  div.white
    v-layout(row, wrap)
      v-flex(xs4)
        search-field(@submit="handleSearch", @clear="handleSearchClear")
      v-flex(v-show="hasDeleteAnyPbehaviorTypeAccess && selected.length", xs4)
        v-btn(@click="showRemoveSelectedPbehaviorTypeModal", icon)
          v-icon delete
    v-data-table(
      v-model="selected",
      :headers="headers",
      :items="pbehaviorTypes",
      :loading="pbehaviorTypesPending",
      :pagination.sync="query",
      item-key="_id",
      select-all,
      expand
    )
      template(slot="items", slot-scope="props")
        tr(@click="props.expanded = !props.expanded")
          td(@click.stop="")
            v-checkbox-functional(v-model="props.selected", primary, hide-details)
          td {{ props.item.name }}
          td
            v-icon {{ props.item.icon_name }}
          td {{ props.item.priority }}
          td
            v-layout
              v-btn.mx-0(
                v-if="hasUpdateAnyPbehaviorTypeAccess",
                icon,
                small,
                @click.stop="showEditPbehaviorTypeModal(props.item)"
              )
                v-icon edit
              v-btn.mx-0(
                v-if="hasDeleteAnyPbehaviorTypeAccess",
                icon,
                small,
                @click.stop="showRemovePbehaviorTypeModal(props.item._id)"
              )
                v-icon(color="error") delete
      template(slot="expand", slot-scope="props")
        pbehavior-types-list-expand-panel(:pbehaviorType="props.item")
</template>

<script>
import { isEqual } from 'lodash';

import { MODALS } from '@/constants';

import entitiesPbehaviorTypesMixin from '@/mixins/entities/pbehavior/types';
import rightsTechnicalPbehaviorTypesMixin from '@/mixins/rights/technical/pbehavior-types';
import pbehaviorTypesQueryMixin from '@/mixins/pbehavior/types/query';

import SearchField from '@/components/forms/fields/search-field.vue';

import PbehaviorTypesListExpandPanel from './partials/pbehavior-types-list-expand-panel.vue';

export default {
  components: {
    SearchField,
    PbehaviorTypesListExpandPanel,
  },
  mixins: [
    entitiesPbehaviorTypesMixin,
    rightsTechnicalPbehaviorTypesMixin,
    pbehaviorTypesQueryMixin,
  ],
  data() {
    return {
      selected: [],
    };
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.icon'),
          value: 'icon_name',
        },
        {
          text: this.$t('common.priority'),
          value: 'priority',
        },
        {
          text: this.$t('common.actionsLabel'),
          sortable: false,
        },
      ];
    },
  },
  watch: {
    query(query, oldQuery) {
      if (!isEqual(query, oldQuery)) {
        this.fetchList();
      }
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      this.fetchPbehaviorTypesList({ params: this.getQuery() });
    },

    showEditPbehaviorTypeModal(pbehaviorType) {
      this.$modals.show({
        name: MODALS.createPbehaviorType,
        config: {
          pbehaviorType,
          action: async (newPbehaviorType) => {
            await this.updatePbehaviorType({
              data: newPbehaviorType,
              id: pbehaviorType._id,
            });
            await this.fetchList();
          },
        },
      });
    },

    showRemovePbehaviorTypeModal(pbehaviorTypeId) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            try {
              await this.removePbehaviorType({ id: pbehaviorTypeId });
              await this.fetchList();
            } catch (err) {
              this.$popups.error({ text: err.error || this.$t('errors.default') });
            }
          },
        },
      });
    },

    showRemoveSelectedPbehaviorTypeModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await Promise.all(this.selected.map(({ _id: id }) => this.removePbehaviorType({ id })));

            await this.fetchList();
            this.selected = [];
          },
        },
      });
    },
  },
};
</script>
