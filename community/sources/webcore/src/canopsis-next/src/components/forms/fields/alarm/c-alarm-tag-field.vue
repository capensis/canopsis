<template>
  <c-lazy-search-field
    v-field="value"
    :items="alarmTags"
    :label="label || $tc('common.tag')"
    :loading="pending"
    :disabled="disabled"
    :name="name"
    :menu-props="{ contentClass: 'c-alarm-tag-field__list' }"
    :has-more="hasMoreTags"
    class="c-alarm-tag-field"
    item-text="value"
    item-value="value"
    hide-details
    multiple
    chips
    dense
    clearable
    autocomplete
    @fetch="fetchTags"
    @fetch:more="fetchMoreTags"
    @update:search="updateSearch"
  >
    <template #selection="{ item, index }">
      <c-alarm-action-chip
        :color="item.color"
        :title="item.value"
        class="c-alarm-tag-field__tag px-2"
        closable
        ellipsis
        @close="removeItemFromArray(index)"
      >
        {{ item.value }}
      </c-alarm-action-chip>
    </template>
    <template #item="{ item, attrs, on, parent }">
      <v-list-item
        class="c-alarm-tag-field__list-item"
        v-bind="attrs"
        v-on="on"
      >
        <v-list-item-action>
          <v-checkbox
            :input-value="attrs.inputValue"
            :color="parent.color"
          />
        </v-list-item-action>
        <v-list-item-content class="c-word-break-all">
          {{ item.value }}
        </v-list-item-content>
      </v-list-item>
    </template>
  </c-lazy-search-field>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { isEmpty, keyBy, pick } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';

import { formArrayMixin } from '@/mixins/form';

const { mapActions: mapAlarmTagActions } = createNamespacedHelpers('alarmTag');

export default {
  mixins: [formArrayMixin],
  props: {
    value: {
      type: [Array],
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'tag',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    limit: {
      type: Number,
      default: PAGINATION_LIMIT,
    },
  },
  data() {
    return {
      tagsByValue: {},
      pending: false,
      pageCount: 1,

      query: {
        page: 1,
        search: null,
      },
    };
  },
  computed: {
    alarmTags() {
      return isEmpty(this.tagsByValue)
        ? this.value.map(value => ({ value }))
        : Object.values(this.tagsByValue);
    },

    hasMoreTags() {
      return this.pageCount > this.query.page;
    },
  },
  mounted() {
    this.fetchTags({
      ...this.getQuery(),
      values: this.value,
    });
  },
  methods: {
    ...mapAlarmTagActions({
      fetchAlarmTagsListWithoutStore: 'fetchListWithoutStore',
    }),

    getQuery() {
      return {
        limit: this.limit,
        page: this.query.page,
        search: this.query.search,
        ...this.params,
      };
    },

    async fetchTags(params = this.getQuery()) {
      try {
        this.pending = true;

        const { data, meta } = await this.fetchAlarmTagsListWithoutStore({
          params,
        });

        this.pageCount = meta.page_count;

        this.tagsByValue = {
          ...(this.query.page !== 1 ? this.tagsByValue : {}),
          ...keyBy(data, 'value'),
          ...pick(this.tagsByValue, this.value),
        };
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },

    fetchMoreTags() {
      this.query.page += 1;

      this.fetchTags();
    },

    updateSearch(search) {
      this.query.search = search;
      this.query.page = 1;

      this.fetchTags();
    },
  },
};
</script>

<style lang="scss">
$selectIconsWidth: 56px;

.c-alarm-tag-field {
  .v-select__selections {
    width: calc(100% - #{$selectIconsWidth});
  }

  &__tag {
    max-width: 100%;
  }

  &__list {
    max-width: 400px;
  }

  &__list-item .v-list-item {
    height: unset !important;
  }
}
</style>
