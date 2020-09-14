<template lang="pug">
  div.white
    v-layout(row, wrap)
      v-flex(xs4)
        search-field(@submit="updateSearchHandler", @clear="clearSearchHandler")
      v-flex(v-show="hasDeleteAnyPbehaviorTypeAccess && selectedTypes.length", xs4)
        v-btn(@click="deleteSelectedTypes", icon)
          v-icon delete
    v-data-table(
      v-model="selected",
      :headers="headers",
      :items="pbehaviorTypes",
      :loading="pending",
      :total-items="totalItems",
      :pagination="pagination",
      item-key="_id",
      select-all,
      expand,
      @update:pagination="$emit('update:pagination', $event)"
    )
      template(slot="items", slot-scope="props")
        tr(@click="props.expanded = !props.expanded")
          td(@click.stop="")
            v-checkbox-functional(
              v-if="props.item.deletable",
              v-model="props.selected",
              primary,
              hide-details
            )
            v-checkbox-functional(v-else, disabled, primary, hide-details)
          td {{ props.item.name }}
          td
            span.pbehavior-type-icon(v-if="props.item.icon_name")
              v-icon(color="white", size="18") {{ props.item.icon_name }}
          td {{ props.item.priority }}
          td
            v-layout
              v-tooltip(bottom, :disabled="props.item.editable")
                v-btn.mx-0(
                  slot="activator",
                  v-if="hasUpdateAnyPbehaviorTypeAccess",
                  :disabled="!props.item.editable",
                  icon,
                  small,
                  @click.stop="$emit('edit', props.item)"
                )
                  v-icon edit
                span {{ $t('pbehaviorTypes.defaultType') }}
              v-tooltip(bottom, :disabled="props.item.deletable")
                v-btn.mx-0(
                  slot="activator",
                  v-if="hasDeleteAnyPbehaviorTypeAccess",
                  :disabled="!props.item.deletable",
                  icon,
                  small,
                  @click.stop="$emit('remove', props.item._id)"
                )
                  v-icon(color="error") delete
                span {{ $t('pbehaviorTypes.usingType') }}
      template(slot="expand", slot-scope="props")
        pbehavior-types-list-expand-panel(:pbehaviorType="props.item")
</template>

<script>
import { omit } from 'lodash';

import rightsTechnicalPbehaviorTypesMixin from '@/mixins/rights/technical/pbehavior-types';

import SearchField from '@/components/forms/fields/search-field.vue';

import PbehaviorTypesListExpandPanel from './partials/pbehavior-types-list-expand-panel.vue';

export default {
  components: {
    SearchField,
    PbehaviorTypesListExpandPanel,
  },
  mixins: [rightsTechnicalPbehaviorTypesMixin],
  props: {
    pbehaviorTypes: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pagination: {
      type: Object,
      required: true,
    },
  },
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
          sortable: false,
        },
        {
          text: this.$t('common.priority'),
          value: 'priority',
        },
        {
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
        },
      ];
    },

    selectedTypes() {
      return this.selected.filter(({ deletable }) => deletable);
    },
  },
  methods: {
    updateSearchHandler(search) {
      this.$emit('update:pagination', { ...this.pagination, search });
    },

    clearSearchHandler() {
      this.$emit('update:pagination', omit(this.pagination, ['search']));
    },

    deleteSelectedTypes() {
      this.$emit('remove-selected', this.selectedTypes);
    },
  },
};
</script>

<style lang="scss" scoped>
  .pbehavior-type-icon {
    display: inline-flex;
    padding: 2px 10px;
    border-radius: 10px;
    background: #17ffff;
  }
  .item-checkbox {
    display: inline-block;
  }
</style>
