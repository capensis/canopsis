<template lang="pug">
  v-layout(column)
    c-entity-field(
      :value="form.entity",
      :label="$t('map.defineEntity')",
      :required="!isLinked",
      :clearable="isLinked",
      :entity-types="entityTypes",
      :item-disabled="isEntityExist",
      :item-text="getItemText",
      return-object,
      autocomplete,
      @input="updateEntity"
    )
      template(v-if="coordinates", #icon="{ item }")
        v-icon.mr-2(v-if="item.coordinates") pin_drop
    c-coordinates-field(v-if="coordinates", v-field="form.coordinates", :disabled="form.is_entity_coordinates")
    c-enabled-field(v-model="isLinked", :label="$t('map.addLink')")
    c-map-field(
      v-show="isLinked",
      v-field="form.map",
      :required="isLinked",
      hide-details
    )
</template>

<script>
import { ENTITY_TYPES } from '@/constants';

import { getMapEntityText } from '@/helpers/map';

import { formMixin } from '@/mixins/form';

export default {
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    coordinates: {
      type: Boolean,
      required: false,
    },
    existsEntities: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    entityTypes() {
      return Object.values(ENTITY_TYPES);
    },

    isLinked: {
      set(value) {
        this.updateField('map', value ? '' : undefined);
      },
      get() {
        return this.form.map !== undefined && this.form.map !== null;
      },
    },
  },
  watch: {
    'form.coordinates': {
      handler(value) {
        this.$emit('fly:coordinates', value);
      },
    },
  },
  methods: {
    updateEntity(entity) {
      if (!this.coordinates) {
        this.updateField('entity', entity?._id ?? '');
        return;
      }

      if (entity) {
        this.updateModel({
          ...this.form,
          entity: entity._id,
          is_entity_coordinates: !!entity.coordinates,
          coordinates: entity.coordinates ?? this.form.coordinates,
        });
      } else {
        this.updateModel({ ...this.form, entity, is_entity_coordinates: false });
      }
    },

    getItemText(item) {
      return getMapEntityText(item);
    },

    isEntityExist(entity) {
      return this.existsEntities.some(id => id === entity._id);
    },
  },
};
</script>
