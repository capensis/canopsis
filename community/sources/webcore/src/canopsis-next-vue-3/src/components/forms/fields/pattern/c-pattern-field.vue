<template>
  <c-select
    :model-value="value"
    :items="itemsWithCustom"
    :label="label || $t('common.pattern')"
    :loading="pending"
    :disabled="disabled"
    :name="name"
    :required="required"
    :return-object="returnObject"
    item-title="title"
    item-value="_id"
    hide-details
    @update:model-value="$emit('update:value', $event)"
  >
    <template #item="{ props: itemProps, item }">
      <v-list-item v-bind="itemProps">
        <template v-if="item.is_corporate" #append>
          <v-icon size="20">
            share
          </v-icon>
        </template>
      </v-list-item>
    </template>
  </c-select>
</template>

<script>
import { MAX_LIMIT, PATTERN_CUSTOM_ITEM_VALUE } from '@/constants';
import { entitiesPatternsMixin } from '@/mixins/entities/pattern';

export default {
  emits: ['update:value'],
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
