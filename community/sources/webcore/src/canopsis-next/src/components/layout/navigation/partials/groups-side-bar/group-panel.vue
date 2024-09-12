<template lang="pug">
  v-expansion-panel-content.secondary.group-item(
    :hide-actions="hideActions",
    :class="{ editing: isEditing }"
  )
    template(#header="")
      div.panel-header
        slot(name="title")
          div.panel-header__title {{ group.title }}
        div.panel-header__actions
          v-icon.mr-2(v-if="group.is_private", small) lock
          v-btn(
            v-if="editable",
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
    editable: {
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

    &__actions {
      flex-shrink: 0;
    }

    &__title {
      width: 100%;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      display: inline-block;
    }
  }
</style>
