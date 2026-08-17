<template>
  <tr v-if="shown" class="alarms-list-table__expand-row">
    <td :colspan="colspan" class="pa-0">
      <transition name="alarms-list-expand-row" @after-leave="hide">
        <div v-show="visible" class="alarms-list-table__expand-row-content">
          <slot />
        </div>
      </transition>
    </td>
  </tr>
</template>

<script>
export default {
  props: {
    expanded: {
      type: Boolean,
      default: false,
    },
    colspan: {
      type: Number,
      required: true,
    },
  },

  data() {
    return {
      shown: this.expanded,
      visible: this.expanded,
    };
  },

  watch: {
    expanded(expanded) {
      if (expanded) {
        this.shown = true;
        this.$nextTick(() => {
          this.visible = true;
        });

        return;
      }

      this.visible = false;
    },
  },

  methods: {
    hide() {
      if (!this.expanded) {
        this.shown = false;
      }
    },
  },
};
</script>

<style lang="scss">
.alarms-list-table__expand-row-content {
  overflow: hidden;
}

.alarms-list-expand-row-enter-active,
.alarms-list-expand-row-leave-active {
  transition:
    opacity 0.15s ease,
    transform 0.15s ease;
}

.alarms-list-expand-row-enter-from,
.alarms-list-expand-row-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
