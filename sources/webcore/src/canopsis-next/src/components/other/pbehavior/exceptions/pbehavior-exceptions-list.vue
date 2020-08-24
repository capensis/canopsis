<template lang="pug">
  div.white
    v-layout(row, wrap)
      v-flex(xs4)
        search-field(@submit="updateSearchHandler", @clear="clearSearchHandler")
      v-flex(v-show="hasDeleteAnyPbehaviorExceptionAccess && selected.length", xs4)
        v-btn(@click="$emit('remove-selected', selected)", icon)
          v-icon delete
    v-data-table(
      v-model="selected",
      :headers="headers",
      :items="pbehaviorExceptions",
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
            v-checkbox-functional(v-model="props.selected", :disabled="!props.item.deletable", primary, hide-details)
          td {{ props.item.name }}
          td
            v-layout
              v-btn.mx-0(
                slot="activator",
                v-if="hasUpdateAnyPbehaviorExceptionAccess",
                icon,
                small,
                @click.stop="$emit('edit', props.item)"
              )
                v-icon edit
              v-tooltip(bottom, :disabled="props.item.deletable")
                v-btn.mx-0(
                  slot="activator",
                  v-if="hasDeleteAnyPbehaviorExceptionAccess",
                  :disabled="!props.item.deletable",
                  icon,
                  small,
                  @click.stop="$emit('remove', props.item._id)"
                )
                  v-icon(color="error") delete
                span {{ $t('pbehaviorExceptions.usingException') }}
      template(slot="expand", slot-scope="props")
        pbehavior-exceptions-list-expand-panel(:pbehavior-exception="props.item")
</template>

<script>
import { omit } from 'lodash';

import rightsTechnicalPbehaviorExceptionsMixin from '@/mixins/rights/technical/pbehavior-exceptions';

import SearchField from '@/components/forms/fields/search-field.vue';

import PbehaviorExceptionsListExpandPanel from './partials/pbehavior-exceptions-list-expand-panel.vue';

export default {
  components: {
    PbehaviorExceptionsListExpandPanel,
    SearchField,
  },
  mixins: [rightsTechnicalPbehaviorExceptionsMixin],
  props: {
    pbehaviorExceptions: {
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
          text: this.$t('common.actionsLabel'),
          sortable: false,
        },
      ];
    },
  },
  methods: {
    updateSearchHandler(search) {
      this.$emit('update:pagination', { ...this.pagination, search });
    },

    clearSearchHandler() {
      this.$emit('update:pagination', omit(this.pagination, ['search']));
    },
  },
};
</script>

<style lang="scss" scoped>
  .item-checkbox {
    display: inline-block;
  }
</style>
