<template lang="pug">
  v-expansion-panel-content.secondary.white--text.group-item(
    :hide-actions="hideActions",
    :class="{ editing: isEditing }"
  )
    template(#header="")
      div.panel-header
        slot(name="title")
          span {{ group.title }}
        v-btn(
          v-show="isEditing",
          :disabled="orderChanged",
          depressed,
          small,
          icon,
          @click.stop="handleChange"
        )
          v-icon(small) edit
    slot
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
