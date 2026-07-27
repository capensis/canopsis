<template>
  <c-layout align="center" wrap>
    <div class="flex-grow-1">
      <v-chip
        v-for="entity in entities"
        :key="getEntityKey(entity)"
        :closable="!isDisabledEntity(entity)"
        :class="{ 'bg-error text-white': isEntityExists(entity) }"
        @click:close="$emit('remove', entity)"
      >
        {{ contentKey ? entity[contentKey] : entity }}
      </v-chip>
    </div>
    <div v-if="clearable && removableEntities.length">
      <c-action-btn
        :tooltip="$t('common.deleteAll')"
        type="delete"
        color="red"
        size="small"
        @click="clear"
      />
    </div>
  </c-layout>
</template>

<script>
export default {
  emits: ['remove', 'clear'],
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
      return this.entities.filter((entity) => this.isDisabledEntity(entity));
    },
  },
  methods: {
    isDisabledEntity(entity) {
      const key = this.getEntityKey(entity);

      return key
        ? this.disabledEntities.some((disabledEntity) => this.getEntityKey(disabledEntity) === key)
        : this.disabledEntities.includes(entity);
    },

    isEntityExists(entity) {
      const key = this.getEntityKey(entity);

      return key
        ? this.existingEntities.some((existEntity) => this.getEntityKey(existEntity) === key)
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
