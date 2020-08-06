<template lang="pug">
  div.white
    v-layout(row, wrap)
      v-flex(xs4)
        search-field(@submit="updateSearchHandler")
      v-flex(v-show="hasDeleteAnyPbehaviorTypeAccess && selected.length", xs4)
        v-btn(@click="$emit('remove-selected', selected)", icon)
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
            v-checkbox-functional(v-model="props.selected", primary, hide-details)
          td {{ props.item.name }}
          td
            span.pbehavior-type-icon
              v-icon(color="white", size="18") {{ props.item.icon_name }}
          td {{ props.item.priority }}
          td
            v-layout
              v-btn.mx-0(
                v-if="hasUpdateAnyPbehaviorTypeAccess",
                icon,
                small,
                @click.stop="$emit('edit', props.item)"
              )
                v-icon edit
              v-btn.mx-0(
                v-if="hasDeleteAnyPbehaviorTypeAccess",
                icon,
                small,
                @click.stop="$emit('remove', props.item._id)"
              )
                v-icon(color="error") delete
      template(slot="expand", slot-scope="props")
        pbehavior-types-list-expand-panel(:pbehaviorType="props.item")
</template>

<script>
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
          sortable: false,
        },
      ];
    },
  },
  methods: {
    updateSearchHandler(search) {
      this.$emit('update:search', search);
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
