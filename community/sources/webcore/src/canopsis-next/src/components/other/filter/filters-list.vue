<template>
  <div>
    <c-alert
      :value="!pending && !filters.length"
      :type="errors.has(name) ? 'error' : 'info'"
    >
      {{ $t('modals.createFilter.emptyFilters') }}
    </c-alert>
    <c-draggable-list-field
      v-field="filters"
      :disabled="!editable"
      handle=".action-drag-handler"
    >
      <filter-tile
        v-for="(filter, index) in filters"
        :key="filter._id"
        :filter="filter"
        :editable="editable"
        @edit="$emit('edit', filter, index)"
        @delete="$emit('delete', filter, index)"
      />
    </c-draggable-list-field>
    <v-btn
      v-if="addable"
      class="v-btn-legacy-m--y"
      color="primary"
      outlined
      @click.prevent="$emit('add', $event)"
    >
      {{ $t('common.addFilter') }}
    </v-btn>
  </div>
</template>

<script>
import { Validator } from 'vee-validate';

import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';
import { validationAttachRequiredMixin } from '@/mixins/form/validation-attach-required';

import FilterTile from './partials/filter-tile.vue';

export default {
  inject: {
    $validator: {
      default: new Validator(),
    },
  },
  components: { FilterTile },
  mixins: [entitiesWidgetMixin, validationAttachRequiredMixin],
  model: {
    prop: 'filters',
    event: 'input',
  },
  props: {
    filters: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: true,
    },
    editable: {
      type: Boolean,
      default: true,
    },
    name: {
      type: String,
      default: 'filters',
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  watch: {
    filters() {
      if (this.required) {
        this.validateRequiredRule();
      }
    },
  },
  mounted() {
    if (this.required) {
      this.attachRequiredRule(() => this.filters.length > 0);
    }
  },
  beforeDestroy() {
    this.detachRequiredRule();
  },
};
</script>
