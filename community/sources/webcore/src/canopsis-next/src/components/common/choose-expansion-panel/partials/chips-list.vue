<template>
  <v-layout
    wrap
    align-center
  >
    <v-flex xs11>
      <v-chip
        v-for="entity in entities"
        :key="getEntityKey(entity)"
        :close="!isDisabledEntity(entity)"
        :class="{ 'error white--text': isEntityExists(entity) }"
        @click:close="$emit('remove', entity)"
      >
        {{ contentKey ? entity[contentKey] : entity }}
      </v-chip>
    </v-flex>
    <v-flex
      v-if="clearable && removableEntities.length"
      xs1
    >
      <c-action-btn
        :tooltip="$t('common.deleteAll')"
        type="delete"
        color="red"
        small
        @click="clear"
      />
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  props: {
    entities: {
      type: Array,
      required: true,
    },
    disabledEntities: {
      type: Array,
      default: () => [],
    },
    existingEntities: {
      type: Array,
      default: () => [],
    },
    itemKey: {
      type: String,
      default: '_id',
    },
    contentKey: {
      type: String,
      required: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    removableEntities() {
      return this.entities.filter(entity => this.isDisabledEntity(entity));
    },
  },
  methods: {
    isDisabledEntity(entity) {
      const key = this.getEntityKey(entity);

      return key
        ? this.disabledEntities.some(disabledEntity => this.getEntityKey(disabledEntity) === key)
        : this.disabledEntities.includes(entity);
    },

    isEntityExists(entity) {
      const key = this.getEntityKey(entity);

      return key
        ? this.existingEntities.some(existEntity => this.getEntityKey(existEntity) === key)
        : this.existingEntities.includes(entity);
    },

    clear() {
      this.$emit('clear', this.removableEntities);
    },

    getEntityKey(entity) {
      return entity[this.itemKey] || entity;
    },
  },
};
</script>
