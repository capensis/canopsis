<template lang="pug">
v-card-text
  v-container
    v-layout(row)
      v-text-field(
        :label="$t('common.name')",
        :value="name",
        @input="(name) => $emit('update:name', name)",
        :error-messages="errors.collect('name')"
        v-validate="'required'",
        data-vv-name="name"
      )
    v-layout(row)
      v-text-field(
        :label="$t('common.description')",
        :value="description",
        @input="(description) => $emit('update:description', description)",
        data-vv-name="description",
        :error-messages="errors.collect('description')"
        multi-line
      )
    v-layout(row)
      v-switch(:label="$t('common.enabled')", :input-value="enabled", @change="(enabled) => $emit('update:enabled', enabled)")
      v-select(
        :items="types",
        :value="entityType",
        item-text="text",
        item-value="value",
        data-vv-name="type",
        v-validate="'required'",
        :error-messages="errors.collect('type')"
        @input="(type) => $emit('update:type', type)",
        label="Type"
        single-line
      )
    v-layout(wrap)
      v-flex(xs12)
        entities-select(label="Impacts", :entities="impact", @updateEntities="updateImpact")
      v-flex(xs12)
        entities-select(label="Dependencies", :entities="depends", @updateEntities="updateDepends")
</template>

<script>
import EntitiesSelect from '@/components/other/context/actions/create-entities/entities-select.vue';

import { MODALS } from '@/constants';

export default {
  name: MODALS.createEntity,
  inject: ['$validator'],
  components: {
    EntitiesSelect,
  },
  props: {
    name: {
      type: String,
      default: '',
    },
    description: {
      type: String,
      default: '',
    },
    type: {
      type: String,
      default: '',
    },
    impact: {
      type: Array,
      default() {
        return [];
      },
    },
    depends: {
      type: Array,
      default() {
        return [];
      },
    },
    enabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      showValidationErrors: true,
      types: [
        {
          text: this.$t('modals.createEntity.fields.types.connector'),
          value: 'connector',
        },
        {
          text: this.$t('modals.createEntity.fields.types.component'),
          value: 'component',
        },
        {
          text: this.$t('modals.createEntity.fields.types.resource'),
          value: 'resource',
        },
      ],
    };
  },
  computed: {
    entityType() {
      let entityType;
      this.types.map((index, item) => {
        if (this.type === item.value) {
          console.log(item);
          return entityType = this.types[index];
        }
        return null;
      });
      return entityType;
    },
  },
  methods: {
    updateImpact(entities) {
      this.$emit('update:impact', entities);
    },
    updateDepends(entities) {
      this.$emit('update:depends', entities);
    },
  },
};
</script>

<style scoped>
  .tooltip {
    flex: 1 1 auto;
  }
</style>
