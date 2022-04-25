<template lang="pug">
  c-select-field(
    v-field="value",
    v-validate="rules",
    :items="itemsWithCustom",
    :label="label || $tc('common.pattern')",
    :loading="pending",
    :disabled="disabled",
    :name="name",
    :return-object="returnObject",
    item-text="title",
    item-value="_id",
    dense,
    hide-details
  )
    template(#item="{ item, tile }")
      v-list-tile(v-bind="tile.props", v-on="tile.on")
        v-list-tile-content {{ item.title }}
        v-icon(v-if="item.is_corporate", size="20") share
</template>

<script>
import { MAX_LIMIT, PATTERN_CUSTOM_ITEM_VALUE } from '@/constants';

import { entitiesPatternsMixin } from '@/mixins/entities/pattern';

export default {
  inject: ['$validator'],
  mixins: [entitiesPatternsMixin],
  props: {
    value: {
      type: [Object, String, Symbol],
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'filter',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    type: {
      type: String,
      required: false,
    },
  },
  data() {
    return {
      items: [],
      pending: false,
    };
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },

    itemsWithCustom() {
      return [
        ...this.items,
        {
          _id: PATTERN_CUSTOM_ITEM_VALUE,
          title: this.$t('common.custom'),
        },
      ];
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    getQuery() {
      const params = { limit: MAX_LIMIT };

      if (this.type) {
        params.type = this.type;
      }

      return params;
    },

    async fetchList() {
      this.pending = true;

      const { data: items } = await this.fetchPatternsListWithoutStore({ params: this.getQuery() });

      this.items = items;
      this.pending = false;
    },
  },
};
</script>
