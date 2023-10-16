<template>
  <v-expansion-panel-content
    class="secondary group-item"
    :hide-actions="hideActions"
    :class="{ editing: isEditing }"
  >
    <template #header="">
      <div class="panel-header">
        <slot name="title">
          <span>{{ group.title }}</span>
        </slot>
        <v-btn
          v-show="isEditing"
          :disabled="orderChanged"
          depressed
          small
          icon
          @click.stop="handleChange"
        >
          <v-icon small>
            edit
          </v-icon>
        </v-btn>
      </div>
    </template>
    <slot />
  </v-expansion-panel-content>
</template>

<script>
export default {
  props: {
    isEditing: {
      type: Boolean,
      default: false,
    },
    group: {
      type: Object,
      required: true,
    },
    orderChanged: {
      type: Boolean,
      default: false,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    handleChange() {
      this.$emit('change');
    },
  },
};
</script>

<style lang="scss" scoped>
  .panel-header {
    max-width: 88%;
    display: flex;
    align-items: center;
    justify-content: space-between;

    span {
      max-width: 100%;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      display: inline-block;
      vertical-align: middle;

      .editing & {
        max-width: 73%;
      }
    }
  }
</style>
